package internal

import (
	"time"
)

type Code struct {
	UserId     string    `json:"userId"`
	Hash       string    `json:"hash"`
	Expiration time.Time `json:"-"`
}

type ClientCode struct {
	UserId string `json:"userId"`
	Code   string `json:"code"`
}
