package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-playground/validator/v10"
	"github.com/kjk/betterguid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/client_golang/prometheus/push"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type (
	Todo struct {
		ID       string   `json:"id"`
		Revision string   `json:"rev"`
		Text     string   `json:"text"`
		Status   string   `json:"status"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}
)

func NewID() string {
	return betterguid.New()
}

func (t Todo) WithNewID(id string) Todo {
	t.ID = id
	return t
}

func (t Todo) WithNewRevision(id string) Todo {
	t.Revision = id
	return t
}

type (
	TodoRepository interface {
		FindByID(id string) (Todo, bool)
		DeleteByID(id string)
		Persist(todo Todo) error
		Update(todo Todo) (Todo, error)
		Prune()
	}

	InMemoryTodoRepository struct {
		mutex       sync.RWMutex
		todoByID    map[string]Todo
		maxCapacity int
	}

	TodoNotFoundError     struct{}
	RevisionMismatchError struct{}
	OutOfCapacityError    struct{}
)

func (e *TodoNotFoundError) Error() string {
	return "Todo not found"
}

func (e *RevisionMismatchError) Error() string {
	return "Revision mismatch"
}

func (e *OutOfCapacityError) Error() string {
	return "Repository is out of capacity"
}

func NewRevisionMismatchError() *RevisionMismatchError {
	return &RevisionMismatchError{}
}

func NewTodoNotFoundError() *TodoNotFoundError {
	return &TodoNotFoundError{}
}

func NewOutOfCapacityError() *OutOfCapacityError {
	return &OutOfCapacityError{}
}

func NewInMemoryTodoRepository(maxCapacity int) TodoRepository {
	return &InMemoryTodoRepository{
		todoByID:    make(map[string]Todo),
		maxCapacity: maxCapacity,
	}
}

func (r *InMemoryTodoRepository) FindByID(id string) (Todo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	t, ok := r.todoByID[id]
	if !ok {
		return Todo{}, false
	}

	return t, true
}

func (r *InMemoryTodoRepository) Persist(todo Todo) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if len(r.todoByID) >= r.maxCapacity {
		return NewOutOfCapacityError()
	}

	r.todoByID[todo.ID] = todo

	return nil
}

func (r *InMemoryTodoRepository) Update(todo Todo) (Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	found, ok := r.todoByID[todo.ID]
	if !ok {
		return Todo{}, NewTodoNotFoundError()
	}

	if found.Revision != todo.Revision {
		return Todo{}, NewRevisionMismatchError()
	}

	r.todoByID[todo.ID] = todo.WithNewRevision(NewID())

	return r.todoByID[todo.ID], nil
}

func (r *InMemoryTodoRepository) Prune() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.todoByID = make(map[string]Todo)
}

func (r *InMemoryTodoRepository) DeleteByID(id string) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	delete(r.todoByID, id)
}

type (
	CreateTodoRequest struct {
		Text     string   `json:"text" validate:"required"`
		Status   string   `json:"status" validate:"required"`
		Category string   `json:"category" validate:"required"`
		Tags     []string `json:"tags" validate:"required"`
	}

	UpdateTodoRequest struct {
		Revision string   `json:"rev" validate:"required"`
		Text     string   `json:"text" validate:"required"`
		Status   string   `json:"status" validate:"required"`
		Category string   `json:"category" validate:"required"`
		Tags     []string `json:"tags" validate:"required"`
	}

	APIError struct {
		Message string   `json:"message"`
		Details []string `json:"details,omitempty"`
	}
)

func (r *CreateTodoRequest) AsNewTodo() Todo {
	return Todo{
		ID:       NewID(),
		Revision: NewID(),
		Text:     r.Text,
		Status:   r.Status,
		Category: r.Category,
		Tags:     r.Tags,
	}
}

func (r *UpdateTodoRequest) AsNewTodo(id string) Todo {
	return Todo{
		ID:       id,
		Revision: r.Revision,
		Text:     r.Text,
		Status:   r.Status,
		Category: r.Category,
		Tags:     r.Tags,
	}
}

func NewAPIInternalServerError() *APIError {
	return &APIError{Message: "Oops"}
}

func NewAPIRevisionMismatchError() *APIError {
	return &APIError{Message: "Revision mismatch"}
}

func NewAPITodoNotFoundError() *APIError {
	return &APIError{Message: "Todo not found"}
}

func NewValidationError(err error) *APIError {
	var details []string
	for _, err := range err.(validator.ValidationErrors) {
		details = append(details, err.Error())
	}
	return &APIError{Message: "Validation failed", Details: details}
}

func NewAPIError(err error) *APIError {
	return &APIError{Message: err.Error()}
}

var (
	ticker = time.NewTicker(2 * time.Second)
	done   = make(chan bool)

	inFlightGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "A gauge of requests currently being served by the wrapped handler.",
	})

	counter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "api_requests_total",
			Help: "A counter for requests to the wrapped handler.",
		},
		[]string{"code", "method"},
	)

	histVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "response_duration_seconds",
			Help:        "A histogram of request latencies.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"handler": "api"},
		},
		[]string{"method"},
	)

	writeHeaderVec = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:        "write_header_duration_seconds",
			Help:        "A histogram of time to first write latencies.",
			Buckets:     prometheus.DefBuckets,
			ConstLabels: prometheus.Labels{"handler": "api"},
		},
		[]string{},
	)

	responseSize = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "push_request_size_bytes",
			Help:    "A histogram of request sizes for requests.",
			Buckets: []float64{200, 500, 900, 1500},
		},
		[]string{},
	)
)

func Instrument(next http.Handler) http.Handler {
	return promhttp.InstrumentHandlerInFlight(inFlightGauge,
		promhttp.InstrumentHandlerCounter(counter,
			promhttp.InstrumentHandlerDuration(histVec,
				promhttp.InstrumentHandlerTimeToWriteHeader(writeHeaderVec,
					promhttp.InstrumentHandlerResponseSize(responseSize, next),
				),
			),
		),
	)
}

func StartPushing(pusher *push.Pusher) {
	for {
		select {
		case <-done:
			return
		case <-ticker.C:
			_ = pusher.Add()
		}
	}
}

type Options struct {
	Cors            bool
	PushgatewayAddr string
}

func NewTodoAPI(repository TodoRepository, options Options) *chi.Mux {
	validate := validator.New()

	r := chi.NewRouter()

	prometheus.MustRegister(inFlightGauge, counter, histVec, writeHeaderVec, responseSize)

	go StartPushing(push.New(options.PushgatewayAddr, "todo").Gatherer(prometheus.DefaultGatherer))

	r.Use(Instrument)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.NoCache)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)

	if options.Cors {
		r.Use(cors.Handler(cors.Options{
			AllowedOrigins:   []string{"http://localhost:9002"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))
	}

	r.Get("/info", Info())
	r.Get("/health", Health())
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/todo", func(r chi.Router) {
		r.Post("/", CreateTodo(repository, validate))
		r.Put("/{id}", UpdateTodoByID(repository, validate))
		r.Get("/{id}", GetTodoByID(repository))
		r.Delete("/{id}", DeleteTodoByID(repository))
	})

	return r
}

func Info() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Info struct {
			Version string `json:"version"`
		}
		respondWithJSON(http.StatusOK, &Info{Version: version}, w)
	}
}

func Health() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Health struct {
			Status string `json:"status"`
		}
		respondWithJSON(http.StatusOK, &Health{Status: "UP"}, w)
	}
}

func CreateTodo(rep TodoRepository, validate *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateTodoRequest

		status, err := Unmarshall(w, r, &request)
		if err != nil {
			respondWithJSON(status, NewAPIError(err), w)
		}

		if err := validate.Struct(&request); err != nil {
			respondWithJSON(http.StatusBadRequest, NewValidationError(err), w)
			return
		}

		todo := request.AsNewTodo()

		if err = rep.Persist(todo); err != nil {
			respondWithJSON(http.StatusTooManyRequests, NewAPIError(err), w)
		}

		respondWithJSON(http.StatusCreated, &todo, w)
	}
}

func GetTodoByID(rep TodoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		t, ok := rep.FindByID(id)
		if !ok {
			respondWithJSON(http.StatusNotFound, NewAPITodoNotFoundError(), w)
			return
		}

		respondWithJSON(http.StatusOK, &t, w)
	}
}

func UpdateTodoByID(rep TodoRepository, validate *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var request UpdateTodoRequest

		status, err := Unmarshall(w, r, &request)
		if err != nil {
			respondWithJSON(status, NewAPIError(err), w)
		}

		if err := validate.Struct(&request); err != nil {
			respondWithJSON(http.StatusBadRequest, NewValidationError(err), w)
			return
		}

		todo, err := rep.Update(request.AsNewTodo(id))
		if err != nil {
			switch err.(type) {
			case *TodoNotFoundError:
				respondWithJSON(http.StatusNotFound, NewAPITodoNotFoundError(), w)
			case *RevisionMismatchError:
				respondWithJSON(http.StatusConflict, NewAPIRevisionMismatchError(), w)
			default:
				respondWithJSON(http.StatusInternalServerError, NewAPIInternalServerError(), w)
			}
			return
		}

		respondWithJSON(http.StatusOK, &todo, w)
	}
}

func DeleteTodoByID(rep TodoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		t, ok := rep.FindByID(id)
		if !ok {
			respondWithJSON(http.StatusNotFound, NewAPITodoNotFoundError(), w)
			return
		}

		rep.DeleteByID(id)

		respondWithJSON(http.StatusOK, &t, w)
	}
}

const (
	// Max request size of 64KB. None of the current API requests for this
	// server are near that limit. Prevents us from unnecessarily parsing JSON
	// payloads that are much large than we anticipate.
	maxBodyBytes = 64_000
)

func Unmarshall(w http.ResponseWriter, r *http.Request, data interface{}) (int, error) {
	if t := r.Header.Get("content-type"); len(t) < 16 || t[:16] != "application/json" {
		return http.StatusUnsupportedMediaType, fmt.Errorf("content-type is not application/json")
	}

	defer r.Body.Close()
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()

	if err := d.Decode(&data); err != nil {
		var syntaxErr *json.SyntaxError
		var unmarshalError *json.UnmarshalTypeError
		switch {
		case errors.As(err, &syntaxErr):
			return http.StatusBadRequest, fmt.Errorf("malformed json at position %v", syntaxErr.Offset)
		case errors.Is(err, io.ErrUnexpectedEOF):
			return http.StatusBadRequest, fmt.Errorf("malformed json")
		case errors.As(err, &unmarshalError):
			return http.StatusBadRequest, fmt.Errorf("invalid value %v at position %v", unmarshalError.Field, unmarshalError.Offset)
		case strings.HasPrefix(err.Error(), "json: unknown field"):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			return http.StatusBadRequest, fmt.Errorf("unknown field %s", fieldName)
		case errors.Is(err, io.EOF):
			return http.StatusBadRequest, fmt.Errorf("body must not be empty")
		case err.Error() == "http: request body too large":
			return http.StatusRequestEntityTooLarge, err
		default:
			return http.StatusInternalServerError, fmt.Errorf("failed to decode json %v", err)
		}
	}
	if d.More() {
		return http.StatusBadRequest, fmt.Errorf("body must contain only one JSON object")
	}

	return http.StatusOK, nil
}

func respondWithJSON(status int, body interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Panic(err)
	}
}

var (
	version = "dev"
	addr    = flag.String("listen-address", ":8080", "The address to listen on for HTTP requests.")
	pgwaddr = flag.String("pushgateway-address", "http://localhost:9091", "The address of the Prometheus Pushgateway.")
	crs     = flag.Bool("cors", false, "Enable CORS headers")
)

func main() {
	flag.Parse()

	repository := NewInMemoryTodoRepository(100)
	mux := NewTodoAPI(repository, Options{
		Cors:            *crs,
		PushgatewayAddr: *pgwaddr,
	})

	fmt.Printf("Server listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
