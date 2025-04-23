package service

import "context"

func (s Service) GetOrderAndUser(ctx context.Context, order string) (string, string, error) {
	return s.Getter.GetOrderAndUser(ctx, order)
}
