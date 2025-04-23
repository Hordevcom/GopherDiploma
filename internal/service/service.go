package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

type OrderUpdater interface {
	UpdateStatus(ctx context.Context, status, orderID, username string) error
	UpdateUserBalance(ctx context.Context, username string, accrual float32, withdraw float32) error
}

type OrderGetter interface {
	GetOrderAndUser(ctx context.Context, orderID string) (order, username string, err error)
}

type OrderAdder interface {
	AddOrderToDB(ctx context.Context, orderID, username string) error
}

type Service struct {
	Updater OrderUpdater
	Getter  OrderGetter
	Adder   OrderAdder
}

func NewService(db storage.PGDB) *Service {
	return &Service{
		Updater: &db,
		Getter:  &db,
		Adder:   &db,
	}
}
