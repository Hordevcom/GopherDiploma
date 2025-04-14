package handlers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/Hordevcom/GopherDiploma/internal/middleware/auth"
)

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

	go pollOrderStatus(string(body), h.Conf.AccurualSystemAddress)

	w.WriteHeader(http.StatusAccepted)
}

type OrderResponce struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func pollOrderStatus(orderNum string, accrual string) {
	url := fmt.Sprintf("%s/api/orders/%s", accrual, orderNum)
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	timeout := time.After(2 * time.Minute)
	attempts := 0
	maxAttempts := 12

	for {
		select {
		case <-ticker.C:
			attempts++
			if attempts > maxAttempts {
				fmt.Println("Превышено количество попыток")
				return
			}
			resp, err := http.Get(url)
			if err != nil {
				fmt.Println("poll error:", err)
				continue
			}

			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()
			if err != nil {
				fmt.Println("read response error:", err)
				continue
			}

			var responce OrderResponce
			fmt.Printf("Order %s — статус: %s, тело: %s\n", orderNum, resp.Status, string(body))
			err = json.Unmarshal(body, &responce)
			if err != nil {
				fmt.Println("Ошибка парсинга:", err)
				return
			}

			fmt.Println("responce: ", responce)

			if resp.Status == "PROCESSED" {
				fmt.Println("Начислено!!!")
				return
			} else {
				fmt.Println(resp.Status, " not equal PROCESSED!")
			}
		case <-timeout:
			fmt.Println("Time is out")
			return
		}
	}
}
