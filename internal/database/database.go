package database

import (
	"crypto/rand"
	"golang.org/x/crypto/bcrypt"
	"notemeal-server/internal"
)

var Db Database

type Database interface {
	CreateOrUpdateCode(id string) (string, error)
	CreateNote(n *internal.Note) error
	CreateToken(userId string, CodeString string) (*internal.ClientToken, error)
	DeleteNote(id string) error
	DeleteUser(id string) error
	GetCode(userId string) (*internal.Code, error)
	GetNote(id string) (*internal.Note, error)
	GetToken(id string) (*internal.Token, error)
	GetUser(id string) (*internal.User, error)
	IsAdmin(userId string) (bool, error)
	IsNoteOwner(noteId string, principalId string) (bool, error)
	Initialize() error
	ListLastModified(userId string) (map[string]int, error)
	UpdateNote(n *internal.Note) error
	SetUser(u *internal.User) error
}

func CreateSecureString() string {
	return rand.Text()
}

func CompareHashAndString(hash string, str string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(str))
}

func HashString(str string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(str), bcrypt.MinCost)
	return string(hash), err
}
