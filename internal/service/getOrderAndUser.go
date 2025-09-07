package service

import "context"

type OrderGetter interface {
	GetOrderAndUser(ctx context.Context, orderID string) (order, username string, err error)
}

func (s Service) GetOrderAndUser(ctx context.Context, order string) (string, string, error) {
	return s.Getter.GetOrderAndUser(ctx, order)
}
