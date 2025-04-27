package internal

import (
	"time"
)

const CodeJsonKey string = "code"

type Code struct {
	UserId     string
	CodeHash   string
	Expiration time.Time
}
