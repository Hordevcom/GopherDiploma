package routes

import (
	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/handlers"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/logging"
	"github.com/Hordevcom/GopherDiploma/internal/service"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
	"github.com/go-chi/chi/v5"
)

func NewRouter(logger logging.Logger, handler handlers.Handler, db *storage.PGDB, conf config.Config, serv service.Service) *chi.Mux {
	router := chi.NewRouter()
	//db handlers.OrderGetter, accrualAddress string, serv service.Service

	router.Use(logger.WithLogging)

	router.Post("/api/user/register", handler.UserRegister)
	router.Post("/api/user/login", handler.UserLogin)
	//router.With(auth.AuthMiddleware).Post("/api/user/orders", handler.OrderLoad)
	router.With(auth.AuthMiddleware).Post("/api/user/orders", handlers.NewOrderLoad(conf.AccurualSystemAddress, serv))
	router.With(auth.AuthMiddleware).Get("/api/user/orders", handler.OrderGet)
	router.With(auth.AuthMiddleware).Get("/api/user/balance", handler.Balance)
	router.With(auth.AuthMiddleware).Post("/api/user/balance/withdraw", handler.BalanceWithdraw)
	router.With(auth.AuthMiddleware).Get("/api/user/withdrawals", handler.Withdraw)

	return router
}
