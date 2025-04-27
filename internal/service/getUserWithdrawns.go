package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
)

type UserWithdrawnsGetter interface {
	GetUserWithdrawns(ctx context.Context, user string) ([]models.UserWithdrawal, error)
}

func (s Service) GetUserWithdrawns(ctx context.Context, user string) ([]models.UserWithdrawal, error) {

	withdrawns, err := s.UserWithdrawnsGetter.GetUserWithdrawns(ctx, user)

	if err != nil {
		return []models.UserWithdrawal{}, err
	}

	return withdrawns, nil

}
