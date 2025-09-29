package server

import (
	"context"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/gin-gonic/gin"
)

type Server struct {
	service Service
	router  *gin.Engine
}

//go:generate mockgen -source=server.go -destination=../mocks/service.go -package=mocks Service
type Service interface {
	CreateUser(ctx context.Context, user *domain.User) error
	UpdateUser(ctx context.Context, user *domain.User) error
	ListUsers(ctx context.Context) ([]*domain.User, error)
	CreateOrder(ctx context.Context, order *domain.Order) error
	UpdateOrder(ctx context.Context, order *domain.Order) error
	GetOrderByID(ctx context.Context, ID string) (*domain.Order, error)

	GetProductByID(ctx context.Context, ID string) (*domain.Product, error)
	ListProducts(ctx context.Context, limit int, offset int) ([]*domain.Product, error)

	AddProductToCategory(ctx context.Context, categoryID string, productID string) error
	RemoveProductFromCategory(ctx context.Context, categoryID string, productID string) error

	GetSupplierByID(ctx context.Context, ID string) (*domain.Supplier, error)
	DeleteSupplierByID(ctx context.Context, ID string) error
}

func (s *Server) Run(addr string) error {
	router := s.SetupRouter()
	return router.Run(addr)

}

func (s *Server) SetupRouter() *gin.Engine {
	s.router = gin.Default()
	// Users
	s.router.POST("/users", s.CreateUserHandler)
	s.router.PUT("/users/:id", s.UpdateUserHandler)
	s.router.GET("/users", s.ListUsersHandler)

	// Products
	s.router.GET("/products/:id", s.GetProductByIDHandler)
	s.router.GET("/products", s.ListProductsHandler)

	// Orders
	s.router.POST("/orders", s.CreateOrderHandler)
	s.router.PUT("/orders/:id", s.UpdateOrderHandler)
	s.router.GET("/orders/:id", s.GetOrderByIDHandler)

	// Categories
	s.router.POST("/categories/:id/products/:productID", s.AddProductToCategoryHandler)
	s.router.DELETE("/categories/:id/products/:productID", s.RemoveProductFromCategoryHandler)

	// Suppliers
	s.router.GET("supplier/:id", s.GetSupplierByIDHandler)
	s.router.DELETE("supplier/:id", s.DeleteSupplierByIDHandler)

	return s.router
}

func NewServer(service Service) *Server {
	return &Server{service: service}
}
