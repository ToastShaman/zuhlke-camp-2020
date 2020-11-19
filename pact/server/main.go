package main

import (
	"encoding/json"
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
	"log"
	"net/http"
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
		Persist(todo Todo)
		Update(todo Todo) (Todo, error)
		Prune()
	}

	InMemoryTodoRepository struct {
		mutex    sync.RWMutex
		todoByID map[string]Todo
	}

	TodoNotFoundError struct {
		message string
	}

	RevisionMismatchError struct {
		message string
	}
)

var (
	todoNotFoundError     = &TodoNotFoundError{message: "Todo not found"}
	revisionMismatchError = &RevisionMismatchError{message: "Revision mismatch"}
	emptyTodo             = new(Todo)
)

func (e *TodoNotFoundError) Error() string {
	return e.message
}

func (e *RevisionMismatchError) Error() string {
	return e.message
}

func NewInMemoryTodoRepository() TodoRepository {
	return &InMemoryTodoRepository{todoByID: make(map[string]Todo)}
}

func (r *InMemoryTodoRepository) FindByID(id string) (Todo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	t, ok := r.todoByID[id]
	if !ok {
		return *emptyTodo, false
	}

	return t, true
}

func (r *InMemoryTodoRepository) Persist(todo Todo) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.todoByID[todo.ID] = todo
}

func (r *InMemoryTodoRepository) Update(todo Todo) (Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	found, ok := r.todoByID[todo.ID]
	if !ok {
		return *emptyTodo, todoNotFoundError
	}

	if found.Revision != todo.Revision {
		return *emptyTodo, revisionMismatchError
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

var (
	APITodoNotFoundError     = &APIError{Message: "Todo not found"}
	APIRevisionMismatchError = &APIError{Message: "Revision mismatch"}
	APIInternalServerError   = &APIError{Message: "Oops"}
)

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

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithJSON(http.StatusInternalServerError, NewAPIError(err), w)
			return
		}

		if err := validate.Struct(&request); err != nil {
			respondWithJSON(http.StatusBadRequest, NewValidationError(err), w)
			return
		}

		todo := request.AsNewTodo()

		rep.Persist(todo)

		respondWithJSON(http.StatusCreated, &todo, w)
	}
}

func GetTodoByID(rep TodoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		t, ok := rep.FindByID(id)
		if !ok {
			respondWithJSON(http.StatusNotFound, &APITodoNotFoundError, w)
			return
		}

		respondWithJSON(http.StatusOK, &t, w)
	}
}

func UpdateTodoByID(rep TodoRepository, validate *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var request UpdateTodoRequest

		if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
			respondWithJSON(http.StatusInternalServerError, NewAPIError(err), w)
			return
		}

		if err := validate.Struct(&request); err != nil {
			respondWithJSON(http.StatusBadRequest, NewValidationError(err), w)
			return
		}

		todo, err := rep.Update(request.AsNewTodo(id))
		if err != nil {
			switch err.(type) {
			case *TodoNotFoundError:
				respondWithJSON(http.StatusNotFound, &APITodoNotFoundError, w)
			case *RevisionMismatchError:
				respondWithJSON(http.StatusConflict, &APIRevisionMismatchError, w)
			default:
				respondWithJSON(http.StatusInternalServerError, &APIInternalServerError, w)
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
			respondWithJSON(http.StatusNotFound, &APITodoNotFoundError, w)
			return
		}

		rep.DeleteByID(id)

		respondWithJSON(http.StatusOK, &t, w)
	}
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
	addr    = flag.String("listen-address", "localhost:3000", "The address to listen on for HTTP requests.")
	pgwaddr = flag.String("pushgateway-address", "http://localhost:9091", "The address of the Prometheus Pushgateway.")
	crs     = flag.Bool("cors", false, "Enable CORS headers")
)

func main() {
	flag.Parse()

	repository := NewInMemoryTodoRepository()
	mux := NewTodoAPI(repository, Options{
		Cors:            *crs,
		PushgatewayAddr: *pgwaddr,
	})

	fmt.Printf("Server listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
