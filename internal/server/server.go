package server

import (
	"context"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/gin-gonic/gin"
)

type Server struct {
	service Service
}

type Service interface {
	CreateUser(ctx context.Context, user *domain.User) error
}

func (s *Server) Run(addr string) error {
	router := gin.Default()

	router.POST("/users", s.CreateUserHandler)

	return router.Run(addr)
}
