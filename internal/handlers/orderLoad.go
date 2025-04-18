package handlers

import (
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/config"
	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/service"
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

func (h *Handler) OrderLoad(w http.ResponseWriter, r *http.Request) {
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
	order, username, err := h.DB.GetOrderAndUser(r.Context(), string(body))

	if err == nil && order == string(body) {

		if user == username {
			http.Error(w, "Order already exist", http.StatusOK)
			return
		} else {
			http.Error(w, "Order upload another user", http.StatusConflict)
			return
		}
	}

	err = h.DB.AddOrderToDB(r.Context(), string(body), user)

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Add to DB error", http.StatusInternalServerError)
		return
	}

	service.PollOrderStatus(r.Context(), string(body), user, h.Conf.AccurualSystemAddress, h.DB)

	w.WriteHeader(http.StatusAccepted)
}
