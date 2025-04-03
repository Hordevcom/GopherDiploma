package handlers

import "github.com/Hordevcom/GopherDiploma/internal/storage"

type Handler struct {
	DB storage.PGDB
}

func NewHandler(DB storage.PGDB) *Handler {
	return &Handler{
		DB: DB,
	}
}
