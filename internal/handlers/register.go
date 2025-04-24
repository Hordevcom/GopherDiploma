package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/service"
	"github.com/go-chi/render"
)

func NewUserRegister(serv service.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if exist := serv.UserChecker.CheckUsernameLogin(r.Context(), user.Username); exist {
			w.WriteHeader(http.StatusConflict)
			return
		}

		err = serv.AddUserToDB(r.Context(), user)

		if err != nil {
			fmt.Printf("err: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
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
		render.JSON(w, r, map[string]string{"message": "User registration complete!"})
	}
}
