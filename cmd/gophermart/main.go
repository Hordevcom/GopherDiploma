package main

import (
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/handlers/routes"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
)

func main() {

	logger := logging.NewLogger()
	router := routes.NewRouter()

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	logger.Infow("Start server", "addr: ", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Fatalw("create server err: ", err)
	}
}
