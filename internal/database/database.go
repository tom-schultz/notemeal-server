package database

import (
	"notemeal-server/internal"
)

var Db Database

type Database interface {
	CreateOrUpdateCode(id string) (string, error)
	CreateToken(userId string, CodeString string) (string, error)
	DeleteNote(id string) error
	DeleteUser(id string) error
	GetNote(id string) (*internal.Note, error)
	GetToken(token string) (*internal.Token, error)
	GetUser(id string) (*internal.User, error)
	IsNoteOwner(noteId string, principalId string) (bool, error)
	Initialize() error
	ListLastModified(userId string) (map[string]int, error)
	UpdateNote(n *internal.Note) error
	SetUser(u *internal.User) error
}
