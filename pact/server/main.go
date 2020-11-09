package main

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-playground/validator/v10"
	"github.com/kjk/betterguid"
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

func (t Todo) NewId() Todo {
	t.Id = betterguid.New()
	return t
}

func (t Todo) NewRevision() Todo {
	t.Revision = betterguid.New()
	return t
}

type (
	TodoRepository interface {
		FindTodoById(id string) (Todo, bool)
		Persist(todo Todo)
		Update(todo Todo) (Todo, error)
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
	emptyTodo             = &Todo{}
)

func (e *TodoNotFoundError) Error() string {
	return e.message
}

func (e *RevisionMismatchError) Error() string {
	return e.message
}

func NewInMemoryTodoRepository() TodoRepository {
	return &InMemoryTodoRepository{todoById: map[string]Todo{}}
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

	t, ok := r.todoById[todo.Id]
	if !ok {
		return *emptyTodo, todoNotFoundError
	}

	if t.Revision != todo.Revision {
		return *emptyTodo, revisionMismatchError
	}

	todoWithNewRevision := todo.NewRevision()

	r.todoById[todo.Id] = todoWithNewRevision

	return todoWithNewRevision, nil
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

func (r *CreateTodoRequest) NewTodo() Todo {
	return Todo{
		Id:       betterguid.New(),
		Revision: betterguid.New(),
		Text:     r.Text,
		Status:   r.Status,
		Category: r.Category,
		Tags:     r.Tags,
	}
}

func (r *UpdateTodoRequest) NewTodo(id string) Todo {
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

func NewTodoAPI(repository TodoRepository) *chi.Mux {
	validate := validator.New()

	r := chi.NewRouter()

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

	r.Route("/todo", func(r chi.Router) {
		r.Post("/", CreateTodo(repository, validate))
		r.Put("/{id}", UpdateTodoById(repository, validate))
		r.Get("/{id}", GetTodoById(repository))
	})

	return r
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

		todo := request.NewTodo()

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

		todo, err := rep.Update(request.NewTodo(id))
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

func main() {
	repository := NewInMemoryTodoRepository()
	mux := NewTodoAPI(repository)
	err := http.ListenAndServe(":3000", mux)
	if err != nil {
		log.Fatal(err)
	}
}
