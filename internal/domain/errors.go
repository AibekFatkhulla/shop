package domain

import "errors"

var (
	ErrorUserNotFound      = errors.New("user not found")
	ErrorUserAlreadyExists = errors.New("user already exists")
	ErrorProductNotFound   = errors.New("product not found")
	ErrorOrderNotFound     = errors.New("order not found")
	ErrorCategoryNotFound  = errors.New("category not found")
	ErrorSupplierNotFound  = errors.New("supplier not found")
)
