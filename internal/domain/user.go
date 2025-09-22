package domain

import "time"

type User struct {
	ID        string
	Name      string
	Email     string
	Password  string
	Number    string
	Address   string
	Balance   int
	CreatedAt time.Time
	UpdatedAt time.Time
}
