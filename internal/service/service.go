package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"time"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/google/uuid"
)

const salt = "6472386&*@^&*@#^&*@#^364732#@&^@*&hjdskdhkjashd38247328&@#*$&@#7283"

type service struct {
	repo Repository
}

type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
}

// CreateUser is a method for creating a new user in the system
func (s *service) CreateUser(ctx context.Context, user *domain.User) error {
	_, err := s.repo.GetByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, domain.ErrorUserNotFound) {
		return err
	}

	if err == nil {
		return domain.ErrorUserAlreadyExists
	}

	user.ID = uuid.New().String()
	now := time.Now()

	user.CreatedAt = now
	user.UpdatedAt = now

	user.Password = hashPassword(user.Password)

	err = s.repo.CreateUser(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// hashPassword is a method for hash users' passwords
func hashPassword(password string) string {
	hasher := sha256.New()

	hasher.Write([]byte(salt + password))

	hashSum := hasher.Sum(nil)

	return string(hashSum)
}
