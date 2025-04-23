package service

import "context"

func (s Service) AddOrderToDB(ctx context.Context, order string, user string) error {
	return s.Adder.AddOrderToDB(ctx, order, user)
}
