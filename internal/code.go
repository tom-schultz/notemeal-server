package internal

import (
	"time"
)

const CodeJsonKey string = "code"

type Code struct {
	UserId     string
	Hash       string
	Expiration time.Time
}
