package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/service"
)

func (h *Handler) Withdraw(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("token")
	user := auth.GetUsername(cookie.Value)

	withdrawals, err := service.GetUserWithdrawns(r.Context(), h.DB, user)

	if err != nil {
		fmt.Println("error: ", err)
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		http.Error(w, "No withdrawals", http.StatusNoContent)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(withdrawals)
	if err != nil {
		http.Error(w, "Problem with encode data", http.StatusInternalServerError)
	}

}
