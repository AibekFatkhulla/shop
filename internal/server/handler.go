package server

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/aibekfatkhulla/shop/internal/domain"
	"github.com/gin-gonic/gin"
)

func (s *Server) CreateUserHandler(c *gin.Context) {
	var user UserDTO

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

func (s *Server) UpdateUserHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
		return
	}
	var user UserDTO
	if err := c.ShouldBind(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	user.ID = id

	if err := s.service.UpdateUser(c.Request.Context(), (*domain.User)(&user)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (s *Server) ListUsersHandler(c *gin.Context) {
	users, err := s.service.ListUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(200, users)
}

func (s *Server) GetProductByIDHandler(c *gin.Context) {
	id := c.Param("id")
	product, err := s.service.GetProductByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
	}
	c.JSON(http.StatusOK, product)
}

func (s *Server) ListProductsHandler(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")
	limit, err := strconv.Atoi(limitStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	offset, err := strconv.Atoi(offsetStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	products, err := s.service.ListProducts(c.Request.Context(), limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
	}
	c.JSON(http.StatusOK, products)
}

func (s *Server) CreateOrderHandler(c *gin.Context) {
	var order domain.Order
	if err := c.ShouldBind(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	now := time.Now()
	order.CreatedAt = now
	if order.Status == "" {
		order.Status = domain.StatusPending
	}
	if err := s.service.CreateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)

}

func (s *Server) UpdateOrderHandler(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing id"})
		return
	}

	var order domain.Order
	if err := c.ShouldBind(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order.ID = id

	if err := s.service.UpdateOrder(c.Request.Context(), &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (s *Server) GetOrderByIDHandler(c *gin.Context) {
	id := c.Param("id")
	order, err := s.service.GetOrderByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

func (s *Server) AddProductToCategoryHandler(c *gin.Context) {
	categoryID := c.Param("id")
	productID := c.Param("productID")

	if err := s.service.AddProductToCategory(c.Request.Context(), categoryID, productID); err != nil {
		switch err {
		case domain.ErrorCategoryNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add product to category"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product added to category"})
}

func (s *Server) RemoveProductFromCategoryHandler(c *gin.Context) {
	categoryID := c.Param("id")
	productID := c.Param("productID")

	if err := s.service.RemoveProductFromCategory(c.Request.Context(), categoryID, productID); err != nil {
		switch err {
		case domain.ErrorCategoryNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		case domain.ErrorProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to remove product from category"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "product removed from category"})
}

func (s *Server) GetSupplierByIDHandler(c *gin.Context) {
	id := c.Param("id")
	supplier, err := s.service.GetSupplierByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	dto := SupplierDTO{
		ID:   supplier.ID,
		Name: supplier.Name,
	}
	c.JSON(http.StatusOK, dto)
}

func (s *Server) DeleteSupplierByIDHandler(c *gin.Context) {
	id := c.Param("id")
	err := s.service.DeleteSupplierByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrorSupplierNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "supplier not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete supplier"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "supplier deleted"})
}
