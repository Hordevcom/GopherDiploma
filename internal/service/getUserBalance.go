package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
)

func (s Service) GetBalanceOfUser(ctx context.Context, user string) (models.UserBalance, error) {
	balance, err := s.GetBalance.GetUserBalance(ctx, user)

	if err != nil {
		return models.UserBalance{}, err
	}

	return balance, nil
}
