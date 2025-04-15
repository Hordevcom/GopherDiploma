package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func GetOrders(ctx context.Context, user string, db storage.PGDB) ([]models.OrderFloat, error) {
	orders, err := db.GetUserOrders(context.Background(), user)
	var ordersFloat []models.OrderFloat

	for _, order := range orders {
		ordersFloat = append(ordersFloat, models.OrderFloat{
			Number:   order.Number,
			Status:   order.Status,
			Accrual:  float64(order.Accrual) / 100,
			UploadAt: order.UploadAt,
		})
	}
	return ordersFloat, err
}
