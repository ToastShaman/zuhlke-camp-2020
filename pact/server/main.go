package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/kjk/betterguid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"sync"
)

type (
	Todo struct {
		Id       string   `json:"id"`
		Revision string   `json:"rev"`
		Text     string   `json:"text"`
		Status   string   `json:"status"`
		Category string   `json:"category"`
		Tags     []string `json:"tags"`
	}
)

func NewId() string {
	return betterguid.New()
}

func (t Todo) WithNewId(id string) Todo {
	t.Id = id
	return t
}

func (t Todo) WithNewRevision(id string) Todo {
	t.Revision = id
	return t
}

type (
	TodoRepository interface {
		FindTodoById(id string) (Todo, bool)
		Persist(todo Todo)
		Update(todo Todo) (Todo, error)
		Prune()
	}

	InMemoryTodoRepository struct {
		mutex    sync.RWMutex
		todoById map[string]Todo
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
	return &InMemoryTodoRepository{todoById: make(map[string]Todo)}
}

func (r *InMemoryTodoRepository) FindTodoById(id string) (Todo, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	t, ok := r.todoById[id]
	if !ok {
		return *emptyTodo, false
	}

	return t, true
}

func (r *InMemoryTodoRepository) Persist(todo Todo) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.todoById[todo.Id] = todo
}

func (r *InMemoryTodoRepository) Update(todo Todo) (Todo, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	found, ok := r.todoById[todo.Id]
	if !ok {
		return *emptyTodo, todoNotFoundError
	}

	if found.Revision != todo.Revision {
		return *emptyTodo, revisionMismatchError
	}

	r.todoById[todo.Id] = todo.WithNewRevision(NewId())

	return r.todoById[todo.Id], nil
}

func (r *InMemoryTodoRepository) Prune() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.todoById = make(map[string]Todo)
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

	ApiError struct {
		Message string   `json:"message"`
		Details []string `json:"details,omitempty"`
	}
)

func (r *CreateTodoRequest) AsNewTodo() Todo {
	return Todo{
		Id:       NewId(),
		Revision: NewId(),
		Text:     r.Text,
		Status:   r.Status,
		Category: r.Category,
		Tags:     r.Tags,
	}
}

func (r *UpdateTodoRequest) AsNewTodo(id string) Todo {
	return Todo{
		Id:       id,
		Revision: r.Revision,
		Text:     r.Text,
		Status:   r.Status,
		Category: r.Category,
		Tags:     r.Tags,
	}
}

var (
	ApiTodoNotFoundError     = &ApiError{Message: "Todo not found"}
	ApiRevisionMismatchError = &ApiError{Message: "Revision mismatch"}
	ApiInternalServerError   = &ApiError{Message: "Oops"}
)

func NewValidationError(err error) *ApiError {
	var details []string
	for _, err := range err.(validator.ValidationErrors) {
		details = append(details, err.Error())
	}
	return &ApiError{Message: "Validation failed", Details: details}
}

func NewApiError(err error) *ApiError {
	return &ApiError{Message: err.Error()}
}

var (
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

func NewTodoAPI(repository TodoRepository) *chi.Mux {
	validate := validator.New()

	r := chi.NewRouter()

	prometheus.MustRegister(inFlightGauge, counter, histVec, writeHeaderVec, responseSize)

	r.Use(func(handler http.Handler) http.Handler {
		return promhttp.InstrumentHandlerInFlight(inFlightGauge,
			promhttp.InstrumentHandlerCounter(counter,
				promhttp.InstrumentHandlerDuration(histVec,
					promhttp.InstrumentHandlerTimeToWriteHeader(writeHeaderVec,
						promhttp.InstrumentHandlerResponseSize(responseSize, handler),
					),
				),
			),
		)
	})
	r.Use(middleware.RequestID)
	r.Use(middleware.NoCache)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Heartbeat("/ping"))
	r.Use(middleware.AllowContentType("application/json"))
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("Content-Type", "application/json")
			next.ServeHTTP(w, r)
		})
	})

	r.Get("/info", Info())
	r.Handle("/metrics", promhttp.Handler())

	r.Route("/todo", func(r chi.Router) {
		r.Post("/", CreateTodo(repository, validate))
		r.Put("/{id}", UpdateTodoById(repository, validate))
		r.Get("/{id}", GetTodoById(repository))
	})

	return r
}

func Info() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		type Info struct {
			Version string `json:"version"`
		}
		respondWith(http.StatusOK, &Info{Version: version}, w)
	}
}

func CreateTodo(rep TodoRepository, validate *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var request CreateTodoRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWith(http.StatusInternalServerError, NewApiError(err), w)
			return
		}

		err = validate.Struct(&request)
		if err != nil {
			respondWith(http.StatusBadRequest, NewValidationError(err), w)
			return
		}

		todo := request.AsNewTodo()

		rep.Persist(todo)

		respondWith(http.StatusCreated, &todo, w)
	}
}

func GetTodoById(rep TodoRepository) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		t, ok := rep.FindTodoById(id)
		if !ok {
			respondWith(http.StatusNotFound, &ApiTodoNotFoundError, w)
			return
		}

		respondWith(http.StatusOK, &t, w)
	}
}

func UpdateTodoById(rep TodoRepository, validate *validator.Validate) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		var request UpdateTodoRequest

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			respondWith(http.StatusInternalServerError, NewApiError(err), w)
			return
		}

		err = validate.Struct(&request)
		if err != nil {
			respondWith(http.StatusBadRequest, NewValidationError(err), w)
			return
		}

		todo, err := rep.Update(request.AsNewTodo(id))
		if err != nil {
			switch err.(type) {
			case *TodoNotFoundError:
				respondWith(http.StatusNotFound, &ApiTodoNotFoundError, w)
				return
			case *RevisionMismatchError:
				respondWith(http.StatusConflict, &ApiRevisionMismatchError, w)
				return
			default:
				respondWith(http.StatusInternalServerError, &ApiInternalServerError, w)
				return
			}
		}

		respondWith(http.StatusOK, &todo, w)
	}
}

func respondWith(status int, body interface{}, w http.ResponseWriter) {
	w.WriteHeader(status)
	err := json.NewEncoder(w).Encode(body)
	if err != nil {
		panic(err)
	}
}

var (
	version = "dev"
	addr    = flag.String("listen-address", "localhost:3000", "The address to listen on for HTTP requests.")
)

func main() {
	flag.Parse()

	repository := NewInMemoryTodoRepository()
	mux := NewTodoAPI(repository)

	fmt.Printf("Server listening on %s\n", *addr)
	log.Fatal(http.ListenAndServe(*addr, mux))
}
