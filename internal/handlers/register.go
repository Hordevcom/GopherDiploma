package handlers

import "net/http"

func (h *Handler) UserRegister(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
}
