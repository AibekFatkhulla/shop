package server

import "time"

type Status string
type UserDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	Number    string    `json:"number"`
	Address   string    `json:"address"`
	Balance   int       `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type OrderDTO struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Status    Status    `json:"status"`
}

type ProductDTO struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Price  int    `json:"price"`
	SKU    string `json:"sku"`
	Amount int    `json:"amount"`
}

type CategoryDTO struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

type SupplierDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
