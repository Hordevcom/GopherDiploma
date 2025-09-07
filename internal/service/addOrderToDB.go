package service

import "context"

type OrderAdder interface {
	AddOrderToDB(ctx context.Context, orderID, username string) error
}

func (s Service) AddOrderToDB(ctx context.Context, order string, user string) error {
	return s.Adder.AddOrderToDB(ctx, order, user)
}
