package service

import (
	"context"
	"fmt"

	"github.com/Hordevcom/GopherDiploma/internal/models"
)

type GetOrder interface {
	GetUserOrders(ctx context.Context, user string) ([]models.Order, error)
}

func (s Service) GetOrders(ctx context.Context, user string) ([]models.OrderFloat, error) {
	orders, err := s.GetOrder.GetUserOrders(ctx, user)

	if err != nil {
		fmt.Println("Error!!!: ", err)
		return nil, err
	}
	userBalance, err := s.GetBalance.GetUserBalance(ctx, user)

	if err != nil {
		fmt.Println("Error!: ", err)
		return nil, err
	}
	var ordersFloat []models.OrderFloat

	for _, order := range orders {
		ordersFloat = append(ordersFloat, models.OrderFloat{
			Number:   order.Number,
			Status:   order.Status,
			Accrual:  float64(userBalance.Current),
			UploadAt: order.UploadAt,
		})
	}
	return ordersFloat, err
}
