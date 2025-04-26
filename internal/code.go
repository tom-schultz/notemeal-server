package internal

import "time"

type Code struct {
	UserId     string
	CodeHash   string
	Expiration time.Time
}
