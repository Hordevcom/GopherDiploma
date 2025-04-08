package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"golang.org/x/crypto/bcrypt"
)

func (h *Handler) UserLogin(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// check user in database
	password := h.DB.GetUserPassword(r.Context(), user.Username)

	if password == "" {
		http.Error(w, "Wrong login or password", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

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
