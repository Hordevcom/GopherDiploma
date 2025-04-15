package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

func GetBalance(ctx context.Context, user string, db storage.PGDB) (models.UserBalance, error) {
	balance, err := db.GetUserBalance(ctx, user)

	if err != nil {
		return models.UserBalance{}, err
	}

	return balance, nil
}
