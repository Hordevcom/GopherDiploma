package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

type OrderUpdater interface {
	UpdateStatus(ctx context.Context, status, orderID, username string) error
	UpdateUserBalance(ctx context.Context, username string, accrual float32, withdraw float32) error
}

type PasswordGetter interface {
	GetUserPassword(ctx context.Context, username string) string
}

type OrderGetter interface {
	GetOrderAndUser(ctx context.Context, orderID string) (order, username string, err error)
}

type OrderAdder interface {
	AddOrderToDB(ctx context.Context, orderID, username string) error
}

type GetOrder interface {
	GetUserOrders(ctx context.Context, user string) ([]models.Order, error)
}

type GetUserBalance interface {
	GetUserBalance(ctx context.Context, user string) (models.UserBalance, error)
}

type CheckUserLogin interface {
	CheckUsernameLogin(ctx context.Context, username string) bool
	AddUserToDB(ctx context.Context, username, password string) error
}

type BalanceUpdater interface {
	UpdateUserBalance(ctx context.Context, user string, accrual, withdrawn float32) error
	SetUserWithdrawn(ctx context.Context, orderNum, user string, withdrawn float32) error
}

type UserWithdrawnsGetter interface {
	GetUserWithdrawns(ctx context.Context, user string) ([]models.UserWithdrawal, error)
}

type Service struct {
	Updater              OrderUpdater
	Getter               OrderGetter
	Adder                OrderAdder
	PassGetter           PasswordGetter
	GetOrder             GetOrder
	GetBalance           GetUserBalance
	UserChecker          CheckUserLogin
	BalanceUpdater       BalanceUpdater
	UserWithdrawnsGetter UserWithdrawnsGetter
}

func NewService(db storage.PGDB) *Service {
	return &Service{
		Updater:              &db,
		Getter:               &db,
		Adder:                &db,
		PassGetter:           &db,
		GetOrder:             &db,
		GetBalance:           &db,
		UserChecker:          &db,
		BalanceUpdater:       &db,
		UserWithdrawnsGetter: &db,
	}
}
