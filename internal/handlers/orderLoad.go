package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/service"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
	"go.uber.org/zap"
)

type Handler struct {
	DB     storage.PGDB
	Conf   config.Config
	logger zap.SugaredLogger
}

func NewHandler(DB storage.PGDB, Conf config.Config, logger zap.SugaredLogger) *Handler {
	return &Handler{
		DB:     DB,
		Conf:   Conf,
		logger: logger,
	}
}

func NewOrderLoad(accrualAddress string, serv service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("%v", err), http.StatusUnprocessableEntity)
			return
		}

		if len(body) == 0 {
			http.Error(w, "url param required", http.StatusBadRequest)
			return
		}

		if !service.LuhnCheck(string(body)) {
			http.Error(w, "Failed Luhn algo", http.StatusUnprocessableEntity)
			return
		}

		cookie, _ := r.Cookie("token")
		user := auth.GetUsername(cookie.Value)

		order, username, err := serv.GetOrderAndUser(r.Context(), string(body))

		if err == nil && order == string(body) {

			if user == username {
				http.Error(w, "Order already exist", http.StatusOK)
				return
			} else {
				http.Error(w, "Order upload another user", http.StatusConflict)
				return
			}
		}

		err = serv.AddOrderToDB(r.Context(), string(body), user)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Add to DB error", http.StatusInternalServerError)
			return
		}

		go serv.PollOrderStatus(r.Context(), string(body), user, accrualAddress)

		w.WriteHeader(http.StatusAccepted)
	}
}
