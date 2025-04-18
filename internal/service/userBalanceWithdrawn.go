package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func BalanceWithdrawn(ctx context.Context, currentBalance float32, withdrawn models.UserWithdrawal, db storage.PGDB, user string) error {

	finalSum := currentBalance - withdrawn.Sum
	err := db.UpdateUserBalance(ctx, user, finalSum, withdrawn.Sum)

	if err != nil {
		return err
	}

	return nil
}
