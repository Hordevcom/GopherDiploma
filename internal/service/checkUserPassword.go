package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"golang.org/x/crypto/bcrypt"
)

type PasswordGetter interface {
	GetUserPassword(ctx context.Context, username string) string
}

func (s Service) CheckUserPassword(ctx context.Context, user models.User) error {
	password := s.PassGetter.GetUserPassword(ctx, user.Username)
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

	return err
}
