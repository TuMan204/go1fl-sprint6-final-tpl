package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

type Server struct {
	logger *log.Logger
	Serv   *http.Server
}

func NewServer(logger *log.Logger) Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.GetHTML)
	mux.HandleFunc("/upload", handlers.PostHTML)

	serv := http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ErrorLog:     logger,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	server := Server{
		logger: logger,
		Serv:   &serv,
	}

	return server
}
