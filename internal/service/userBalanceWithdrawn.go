package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
)

type BalanceUpdater interface {
	UpdateUserBalance(ctx context.Context, user string, accrual, withdrawn float32) error
	SetUserWithdrawn(ctx context.Context, orderNum, user string, withdrawn float32) error
}

func (s Service) BalanceWithdrawn(ctx context.Context, currentBalance float32, withdrawn models.UserWithdrawal, user string) error {
	finalSum := currentBalance - withdrawn.Sum

	err := s.BalanceUpdater.UpdateUserBalance(ctx, user, finalSum, withdrawn.Sum)

	if err != nil {
		return err
	}

	err = s.BalanceUpdater.SetUserWithdrawn(ctx, withdrawn.OrderNum, user, withdrawn.Sum)

	if err != nil {
		return err
	}

	return nil
}
