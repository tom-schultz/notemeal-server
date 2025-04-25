package main

import "time"

type TokenCode struct {
	UserId     string
	CodeHash   string
	Expiration time.Time
}
