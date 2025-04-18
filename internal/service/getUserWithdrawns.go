package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func GetUserWithdrawns(ctx context.Context, db storage.PGDB, user string) ([]models.UserWithdrawal, error) {

	withdrawns, err := db.GetUserWithdrawns(ctx, user)

	if err != nil {
		return []models.UserWithdrawal{}, err
	}

	return withdrawns, nil

}
