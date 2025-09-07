package server

import (
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/go-chi/chi/v5"
)

func NewServer(conf config.Config, router *chi.Mux) http.Server {
	return http.Server{
		Addr:    conf.ServerAdress,
		Handler: router,
	}
}
