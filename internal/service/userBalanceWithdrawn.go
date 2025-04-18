package service

import (
	"context"
	"fmt"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func BalanceWithdrawn(ctx context.Context, currentBalance float32, withdrawn models.UserWithdrawal, db storage.PGDB, user string) error {
	fmt.Println("currentBalance is: ", withdrawn.Sum)
	fmt.Println("Sum is: ", withdrawn.Sum)
	finalSum := currentBalance - withdrawn.Sum
	fmt.Println("finalSum is: ", finalSum)
	err := db.UpdateUserBalance(ctx, user, finalSum, withdrawn.Sum)

	if err != nil {
		return err
	}

	err = db.SetUserWithdrawn(ctx, withdrawn.OrderNum, user, withdrawn.Sum)

	if err != nil {
		return err
	}

	return nil
}
