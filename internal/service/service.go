package service

import "context"

type OrderUpdater interface {
	UpdateStatus(ctx context.Context, status, orderID, username string) error
	UpdateUserBalance(ctx context.Context, username string, accrual float32, withdraw float32) error
}

type Service struct {
	Updater OrderUpdater
}

func NewService(updater OrderUpdater) *Service {
	return &Service{
		Updater: updater,
	}
}
