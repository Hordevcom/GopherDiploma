package main

import (
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/handlers"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/Hordevcom/GopherDiploma/internal/routes"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func main() {

	logger := logging.NewLogger()
	handler := handlers.NewHandler()
	router := routes.NewRouter(*logger, *handler)

	server := &http.Server{
		Addr:    "localhost:8080",
		Handler: router,
	}

	storage.InitMigrations(logger.Logger)

	logger.Logger.Infow("Start server", "addr: ", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatalw("create server err: ", err)
	}
}
