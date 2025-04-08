package routes

import (
	"github.com/Hordevcom/GopherDiploma/internal/handlers"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/go-chi/chi/v5"
)

func NewRouter(logger logging.Logger, handler handlers.Handler) *chi.Mux {
	router := chi.NewRouter()

	router.Use(logger.WithLogging)

	router.Post("/api/user/register", handler.UserRegister)
	router.Post("/api/user/login", handler.UserLogin)
	router.With(auth.AuthMiddleware).Post("/api/user/orders", handler.OrderLoad)

	return router
}
