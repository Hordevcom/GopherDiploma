package handlers

import (
	"net/http"
)

func (h *Handler) OrderGet(w http.ResponseWriter, r *http.Request) {
	// cookie, _ := r.Cookie("token")
	// user := auth.GetUsername(cookie.Value)

	// orderResp, err := service.GetOrders(r.Context(), user, h.DB)
	// if err != nil {
	// 	fmt.Println(err)
	// 	http.Error(w, "Problem with data", http.StatusInternalServerError)
	// }

	// if len(orderResp) == 0 {
	// 	w.Header().Set("Content-Type", "application/json")
	// 	w.WriteHeader(http.StatusNoContent)
	// 	return
	// }

	// w.Header().Set("Content-Type", "application/json")

	// err = json.NewEncoder(w).Encode(orderResp)
	// if err != nil {
	// 	http.Error(w, "Problem with encode data", http.StatusInternalServerError)
	// }
}
