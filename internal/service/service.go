package service

import (
	"github.com/Hordevcom/GopherDiploma/internal/storage"
)

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
