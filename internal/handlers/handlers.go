package handlers

import (
	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

type Handler struct {
	DB   storage.PGDB
	Conf config.Config
}

func NewHandler(DB storage.PGDB, Conf config.Config) *Handler {
	return &Handler{
		DB:   DB,
		Conf: Conf,
	}
}
