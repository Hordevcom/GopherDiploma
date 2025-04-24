package main

import (
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/Hordevcom/GopherDiploma/internal/routes"
	"github.com/Hordevcom/GopherDiploma/internal/service"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func main() {

	logger := logging.NewLogger()
	conf := config.NewConfig()
	DB := storage.NewPGDB(conf, *logger)
	services := service.NewService(*DB)
	router := routes.NewRouter(*logger, DB, conf, *services)

	server := &http.Server{
		Addr:    conf.ServerAdress,
		Handler: router,
	}

	logger.Logger.Infow("Start server", "addr: ", server.Addr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Logger.Fatalw("create server err: ", err)
	}
}
