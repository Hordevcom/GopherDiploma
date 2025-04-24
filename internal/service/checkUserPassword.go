package service

import (
	"context"

	"github.com/Hordevcom/GopherDiploma/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (s Service) CheckUserPassword(ctx context.Context, user models.User) error {
	password := s.PassGetter.GetUserPassword(ctx, user.Username)
	err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password))

	return err
}
