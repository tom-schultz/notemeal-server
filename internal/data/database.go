package data

import (
	"notemeal-server/internal"
)

type Datasource interface {
	DeleteCode(userId string) error
	DeleteNote(id string) error
	DeleteUser(id string) error
	GetCode(userId string) (*internal.Code, error)
	GetNote(id string) (*internal.Note, error)
	GetNotesByUser(userId string) ([]*internal.Note, error)
	GetToken(id string) (*internal.Token, error)
	GetUser(id string) (*internal.User, error)
	UpdateCode(code *internal.Code) error
	UpdateNote(note *internal.Note) error
	UpdateToken(token *internal.Token) error
	UpdateUser(user *internal.User) error
}
