package apiserver

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/gorilla/mux"
	"github.com/jacobfire/http-rest-api/store"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"time"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
	store *store.Store
}

func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	if err := s.configureStore(); err != nil {
		return err
	}
	s.logger.Info("Starting server...")

	server := s.configureServer()
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
	return nil
}

// Configure Store
func (s *APIServer) configureStore() error {
	st := store.New(s.config.Store)

	log.Println("Configuring store")
	if err := st.Open(); err != nil {
		return err
	}
	s.store = st

	return nil
}

// Handle Hello function
func (s * APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)
	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureServer() http.Server {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/analyse", func(writer http.ResponseWriter, request *http.Request) {
		writer.Write([]byte("Analyse"))
	})
	s.router.HandleFunc("/catalog/{category}/{id:[0-9]+}", s.categoryHandler())
	server := http.Server {
		Addr: s.config.BindAddr,
		Handler: s.router,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	return server
}

func (s * APIServer) categoryHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		io.WriteString(w, "Article")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Category: %v ID:%v \n", vars["category"], vars["id"])
	}
}

func (s *APIServer) Migrate() error {
	log.Println(s.config.Store.DatabaseURL)
	m, err := migrate.New(
		"/migrations",
		s.config.Store.DatabaseURL,
		)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		return err
	}

	return nil
}