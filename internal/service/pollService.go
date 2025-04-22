package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type OrderResponce struct {
	Order   string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual"`
}

func (s Service) PollOrderStatus(ctx context.Context, orderNum, user string, accrual string) {
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
			err = json.Unmarshal(body, &responce)
			if err != nil {
				fmt.Println("Ошибка парсинга:", err)
				return
			}

			if responce.Status == "PROCESSED" {
				err := s.Updater.UpdateStatus(ctx, responce.Status, responce.Order, user)

				if err != nil {
					fmt.Println("Error in update db: ", err)
					return
				}

				err = s.Updater.UpdateUserBalance(ctx, user, float32(responce.Accrual), 0)

				if err != nil {
					fmt.Println("Error in update db: ", err)
					return
				}
				return
			} else {
				fmt.Println(responce.Status, " not equal PROCESSED!")
			}
		case <-timeout:
			fmt.Println("Time is out")
			return
		}
	}
}
