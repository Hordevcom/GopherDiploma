package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/service"
)

func (h *Handler) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("token")
	user := auth.GetUsername(cookie.Value)

	// get data from json POST
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("%v", err), http.StatusUnprocessableEntity)
		return
	}

	if len(body) == 0 {
		http.Error(w, "url param required", http.StatusBadRequest)
		return
	}

	var data models.UserWithdrawal
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
	}

	// get user balance()
	userBalance, err := service.GetBalance(r.Context(), user, h.DB)

	if err != nil {
		fmt.Println("error: ", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	// check balances and if false = return 402
	if userBalance.Current < data.Sum {
		http.Error(w, "Not enough balance", http.StatusPaymentRequired)
		return
	}

	err = service.BalanceWithdrawn(r.Context(), userBalance.Current, data, h.DB, user)

	if err != nil {
		fmt.Println("error: ", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
