package model

import (
	"crypto/rand"
	"log/slog"
	"notemeal-server/internal"
	"notemeal-server/internal/data"
	"time"
)

type Model struct {
	ds data.Datasource
}

func NewModel(ds data.Datasource) Model {
	return Model{ds: ds}
}

func (m *Model) CreateOrUpdateCode(userId string) (string, error) {
	expiration := time.Now().Add(time.Hour)
	codeStr := internal.CreateSecureString()
	hash, err := internal.HashString(codeStr)

	if err != nil {
		return "", err
	}

	code := &internal.Code{
		UserId:     userId,
		Hash:       hash,
		Expiration: expiration,
	}

	err = m.ds.UpdateCode(code)

	if err != nil {
		return "", err
	}

	return codeStr, nil
}

func (m *Model) CreateNote(newNote *internal.Note) error {
	note, err := m.ds.GetNote(newNote.Id)

	if err != nil {
		return err
	}

	if note != nil {
		return internal.Error{"Note already exists!"}
	}

	err = m.ds.UpdateNote(newNote)

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) CreateToken(userId string, codeString string) (*internal.ClientToken, error) {
	code, err := m.ds.GetCode(userId)

	if err != nil {
		return nil, err
	}

	if code == nil {
		return nil, nil
	}

	if !internal.CompareHashAndString(code.Hash, codeString) {
		return nil, nil
	}

	if code.Expiration.Before(time.Now()) {
		slog.Error("Cannot create tokenStr with expired code!", "user", userId)
		return nil, nil
	}

	tokenStr := internal.CreateSecureString()
	tokenHash, err := internal.HashString(tokenStr)

	if err != nil {
		return nil, err
	}
	id := rand.Text()

	tokenId, err := m.ds.GetToken(id)

	for tokenId != nil {
		id = rand.Text()
		tokenId, err = m.ds.GetToken(id)
	}

	if err != nil {
		return nil, err
	}

	token := &internal.Token{
		Id:     id,
		Hash:   tokenHash,
		UserId: userId}

	clientToken := &internal.ClientToken{
		Id:    id,
		Token: tokenStr}

	err = m.ds.UpdateToken(token)

	if err != nil {
		return nil, err
	}

	err = m.ds.DeleteCode(userId)

	if err != nil {
		return nil, err
	}

	return clientToken, nil
}

func (m *Model) DeleteUser(id string) error {
	err := m.ds.DeleteUser(id)

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) DeleteNote(id string) error {
	err := m.ds.DeleteNote(id)

	if err != nil {
		return err
	}

	return nil
}

func (m *Model) GetCode(userId string) (*internal.Code, error) {
	code, err := m.ds.GetCode(userId)

	if err != nil {
		return nil, err
	}

	return code, nil
}

func (m *Model) GetNote(id string) (*internal.Note, error) {
	note, err := m.ds.GetNote(id)

	if err != nil {
		return nil, err
	}

	return note, nil
}

func (m *Model) GetToken(tokenId string) (*internal.Token, error) {
	token, err := m.ds.GetToken(tokenId)

	if err != nil {
		return nil, err
	}

	return token, nil
}

func (m *Model) GetUser(id string) (*internal.User, error) {
	user, err := m.ds.GetUser(id)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (m *Model) IsAdmin(userId string) (bool, error) {
	user, err := m.ds.GetUser(userId)

	if err != nil {
		return false, err
	}

	return user.IsAdmin, nil
}

func (m *Model) IsNoteOwner(noteId string, principalId string) (bool, error) {
	note, err := m.ds.GetNote(noteId)

	if err != nil {
		return false, err
	}

	if note == nil {
		return false, internal.Error{"note does not exist, cannot be owner!!"}
	}

	return note.UserId == principalId, nil
}

func (m *Model) ListLastModified(userId string) (map[string]int, error) {
	data := make(map[string]int)
	notes, err := m.ds.GetNotesByUser(userId)

	if err != nil {
		return nil, err
	}

	for _, note := range notes {
		if note.UserId == userId {
			data[note.Id] = note.LastModified
		}
	}

	return data, nil
}

func (m *Model) UpdateNote(newNote *internal.Note) error {
	oldNote, err := m.ds.GetNote(newNote.Id)

	if err != nil {
		return err
	}

	if oldNote == nil {
		return internal.Error{"Could not find oldNote: " + newNote.Id}
	}

	oldNote.LastModified = newNote.LastModified
	oldNote.Text = newNote.Text
	oldNote.Title = newNote.Title

	return m.ds.UpdateNote(oldNote)
}

func (m *Model) SetUser(u *internal.User) error {
	user, err := m.ds.GetUser(u.Id)

	if err != nil {
		return err
	}

	if user == nil {
		return internal.Error{"user not found in database!"}
	}

	user.Email = u.Email
	return m.ds.UpdateUser(user)
}
