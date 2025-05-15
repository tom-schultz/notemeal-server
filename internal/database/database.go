package database

import (
	"notemeal-server/internal"
)

type Database interface {
	DeleteCode(userId string) error
	DeleteNote(id string) error
	DeleteUser(id string) error
	GetCode(userId string) (*internal.Code, error)
	GetNote(id string) (*internal.Note, error)
	GetNotesByUser(userId string) ([]*internal.Note, error)
	GetToken(id string) (*internal.Token, error)
	GetUser(id string) (*internal.User, error)
	initialize() error
	StoreCode(code *internal.Code) error
	StoreNote(note *internal.Note) error
	StoreToken(token *internal.Token) error
}
