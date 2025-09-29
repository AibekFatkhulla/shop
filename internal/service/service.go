package service

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/aibekfatkhulla/shop/internal/server"
	"github.com/google/uuid"
)

const salt = "6472386&*@^&*@#^&*@#^364732#@&^@*&hjdskdhkjashd38247328&@#*$&@#7283"

type service struct {
	repo Repository
}

//go:generate mockgen -source=service.go -destination=../mocks/repository.go -package=mocks Repository
type Repository interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	GetUserByID(ctx context.Context, id string) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	ListUsers(ctx context.Context) ([]*domain.User, error)

	GetProductByID(ctx context.Context, id string) (*domain.Product, error)
	ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error)

	CreateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, ID string) (*domain.Order, error)
	UpdateOrder(ctx context.Context, order *domain.Order) error

	AddProductToCategory(ctx context.Context, categoryID, productID string) error
	RemoveProductFromCategory(ctx context.Context, categoryID, productID string) error

	GetSupplierByID(ctx context.Context, id string) (*domain.Supplier, error)
	DeleteSupplierByID(ctx context.Context, id string) error
}

func NewService(repo Repository) server.Service {
	return &service{repo: repo}
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

	return fmt.Sprintf("%x", hashSum)
}

// UpdateUser is a method for updating user
func (s *service) UpdateUser(ctx context.Context, user *domain.User) error {
	existingUser, err := s.repo.GetUserByID(ctx, user.ID)

	if err != nil {
		return err
	}

	user.Balance = existingUser.Balance

	if user.Password == "" {
		user.Password = existingUser.Password
	} else {
		user.Password = hashPassword(user.Password)
	}

	now := time.Now()
	user.UpdatedAt = now

	return s.repo.UpdateUser(ctx, user)
}

func (s *service) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	user, err := s.repo.GetUserByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *service) ListUsers(ctx context.Context) ([]*domain.User, error) {
	return s.repo.ListUsers(ctx)
}

// GetProductByID is a method for finding a product by products ID
func (s *service) GetProductByID(ctx context.Context, id string) (*domain.Product, error) {
	product, err := s.repo.GetProductByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *service) ListProducts(ctx context.Context, limit, offset int) ([]*domain.Product, error) {
	if limit < 1 {
		limit = 10
	}
	if offset < 0 {
		offset = 0
	}

	products, err := s.repo.ListProducts(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *service) CreateOrder(ctx context.Context, order *domain.Order) error {
	now := time.Now()
	order.CreatedAt = now
	order.UpdatedAt = now

	if order.Status == "" {
		order.Status = domain.StatusPending
	}

	return s.repo.CreateOrder(ctx, order)
}

func (s *service) GetOrderByID(ctx context.Context, id string) (*domain.Order, error) {
	return s.repo.GetOrderByID(ctx, id)
}

func (s *service) UpdateOrder(ctx context.Context, order *domain.Order) error {
	if order == nil {
		return domain.ErrorOrderNotFound
	}

	_, err := s.repo.GetOrderByID(ctx, order.ID)
	if err != nil {
		return err
	}

	if order.Status == "" {
		order.Status = domain.StatusCompleted
	}
	now := time.Now()
	order.UpdatedAt = now

	if err := s.repo.UpdateOrder(ctx, order); err != nil {
		return err
	}
	return nil
}

func (s *service) AddProductToCategory(ctx context.Context, categoryID, productID string) error {
	return s.repo.AddProductToCategory(ctx, categoryID, productID)
}

func (s *service) RemoveProductFromCategory(ctx context.Context, categoryID, productID string) error {
	return s.repo.RemoveProductFromCategory(ctx, categoryID, productID)
}

func (s *service) GetSupplierByID(ctx context.Context, id string) (*domain.Supplier, error) {
	return s.repo.GetSupplierByID(ctx, id)
}

func (s *service) DeleteSupplierByID(ctx context.Context, id string) error {
	return s.repo.DeleteSupplierByID(ctx, id)
}
