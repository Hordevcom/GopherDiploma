package service

import (
	"context"
	"fmt"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func GetOrders(ctx context.Context, user string, db storage.PGDB) ([]models.OrderFloat, error) {
	orders, err := db.GetUserOrders(ctx, user)

	if err != nil {
		fmt.Println("Error!!!: ", err)
		return nil, err
	}
	userBalance, err := db.GetUserBalance(ctx, user)

	if err != nil {
		fmt.Println("Error!!!: ", err)
		return nil, err
	}
	var ordersFloat []models.OrderFloat

	for _, order := range orders {
		ordersFloat = append(ordersFloat, models.OrderFloat{
			Number:   order.Number,
			Status:   order.Status,
			Accrual:  float64(userBalance.Current) / 100,
			UploadAt: order.UploadAt,
		})
	}
	return ordersFloat, err
}
