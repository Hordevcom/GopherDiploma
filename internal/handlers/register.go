package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/go-chi/render"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password,omitempty"`
}

func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {
	// get info with unmarshal (login pass)
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check login in DB if already exist
	if exist := h.DB.CheckUsernameLogin(r.Context(), user.Username); exist {
		w.WriteHeader(http.StatusConflict)
		return
	}

	// hashing password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// add user in db (username, hashed password)
	err = h.DB.AddUserToDB(r.Context(), user.Username, string(hashedPassword))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, _ := auth.BuildJWTString()
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	render.JSON(w, r, map[string]string{"message": "User registration complete!"})
}
