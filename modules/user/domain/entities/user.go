package entities

import "time"

type User struct {
	ID        int32
	Email     string
	Password  string
	FullName  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
