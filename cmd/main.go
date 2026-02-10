package main

import (
	"log"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	var logger log.Logger
	server := server.NewServer(&logger)

	err := server.Serv.ListenAndServe()
	if err != nil {
		logger.Fatal()
	}
}
