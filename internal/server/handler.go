package server

import (
	"net/http"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateUserHandler(c *gin.Context) {
	var user userDTO

	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if user.Email == "" || user.Password == "" || user.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required fields"})
		return
	}

	err := s.service.CreateUser(c, &domain.User{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Password:  user.Password,
		Number:    user.Number,
		Address:   user.Address,
		Balance:   user.Balance,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	})
	if err != nil {
		if err == domain.ErrorUserAlreadyExists {
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
}
