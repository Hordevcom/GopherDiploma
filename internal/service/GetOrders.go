package service

import (
	"context"
	"fmt"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func GetOrders(ctx context.Context, user string, db storage.PGDB) ([]models.Order, error) {
	orders, err := db.GetUserOrders(context.Background(), user)
	fmt.Println(orders)
	return orders, err
}
