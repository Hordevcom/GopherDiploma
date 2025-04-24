package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) AddUserToDB(ctx context.Context, user models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = s.UserChecker.AddUserToDB(ctx, user.Username, string(hashedPassword))

	if err != nil {
		return err
	}

	return nil
}
