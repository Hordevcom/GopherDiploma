package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
)

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
