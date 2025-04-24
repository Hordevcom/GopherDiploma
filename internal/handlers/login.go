package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/service"
)

func NewUserLogin(serv service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		// check user in database
		err = serv.CheckUserPassword(r.Context(), user)

		if err != nil {
			http.Error(w, "Wrong login or password", http.StatusUnauthorized)
			return
		}

		token, _ := auth.BuildJWTString(user.Username)
		cookie := &http.Cookie{
			Name:     "token",
			Value:    token,
			HttpOnly: true,
		}
		http.SetCookie(w, cookie)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
