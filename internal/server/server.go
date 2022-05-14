package server

import (
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"rag/internal"
	"rag/internal/models"
	"time"
)

type Server struct {
	*http.Server
	logger      *zap.Logger
	userService UserServicier
	settings    internal.Settings
}

type UserServicier interface {
	CreateUser(request CreateUserRequest) (*int, error)
	UpdateUser(request UpdateUserRequest) (*models.User, error)
	DeleteUser(userID int) error
	GetAllUsers() ([]*models.User, error)
	GetUser(id int) (*models.User, error)
}

func NewServer(us UserServicier, logger *zap.Logger, s *internal.Settings) *Server {
	srv := &Server{
		userService: us,
		logger:      logger,
		settings:    *s,
	}

	srv.Server = &http.Server{
		Addr:           s.Database.PORT,
		Handler:        srv.Handler(),
		MaxHeaderBytes: 1 << 20,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
	}
	return srv
}

func (s *Server) Handler() *mux.Router {
	r := mux.NewRouter()

	router := r.PathPrefix("/rag/v1").Subrouter()
	// check server status
	router.HandleFunc("/check", s.HealthCheckHandler).Methods(http.MethodGet)
	// create a user
	router.HandleFunc("/users/create", s.CreateUser).Methods(http.MethodPost)
	// update a user
	router.HandleFunc("/users/update", s.UpdateUser).Methods(http.MethodPut)
	// delete a user
	router.HandleFunc("/users/delete/{id}", s.DeleteUser).Methods(http.MethodDelete)
	// get all users
	router.HandleFunc("/users/all", s.GetAllUsers).Methods(http.MethodGet)
	// get a user
	router.HandleFunc("/users/{id}", s.GetUser).Methods(http.MethodGet)
	return r
}
