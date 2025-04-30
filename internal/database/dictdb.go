package database

import (
	"crypto/rand"
	"log"
	"log/slog"
	"maps"
	"notemeal-server/internal"
	"slices"
	"time"
)

type dictDb struct {
	initialized bool
	admins      []string
	notes       map[string]*internal.Note
	tokens      map[string]*internal.Token
	codes       map[string]*internal.Code
	users       map[string]*internal.User
}

func DictDb() {
	Db = &dictDb{}

	if err := Db.Initialize(); err != nil {
		log.Fatal(err)
	}
}

func (db *dictDb) CreateOrUpdateCode(userId string) (string, error) {
	if !db.initialized {
		return "", DbError{"Database not initialized!"}
	}

	expiration := time.Now().Add(time.Hour)
	code := CreateSecureString()
	hash, err := HashString(code)

	if err != nil {
		return "", err
	}

	db.codes[userId] = &internal.Code{
		UserId:     userId,
		Hash:       hash,
		Expiration: expiration,
	}

	return code, nil
}

func (db *dictDb) CreateToken(userId string, codeString string) (*internal.ClientToken, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	code := db.codes[userId]

	if code == nil {
		return nil, nil
	}

	err := CompareHashAndString(code.Hash, codeString)

	if err != nil {
		return nil, nil
	}

	if code.Expiration.Before(time.Now()) {
		slog.Error("Cannot create tokenStr with expired code!", "user", userId)
		return nil, nil
	}

	tokenStr := CreateSecureString()
	tokenHash, err := HashString(tokenStr)

	if err != nil {
		return nil, err
	}

	id := rand.Text()
	_, ok := db.tokens[id]

	for ok {
		id = rand.Text()
		_, ok = db.tokens[id]
	}

	db.tokens[id] = &internal.Token{
		Id:     id,
		Hash:   tokenHash,
		UserId: userId}

	clientToken := &internal.ClientToken{
		Id:    id,
		Token: tokenStr}

	delete(db.codes, userId)
	return clientToken, nil
}

func (db *dictDb) DeleteUser(id string) error {
	if !db.initialized {
		return DbError{"Database not initialized!"}
	}

	delete(db.users, id)
	return nil
}

func (db *dictDb) DeleteNote(id string) error {
	if !db.initialized {
		return DbError{"Database not initialized!"}
	}

	delete(db.notes, id)
	return nil
}

func (db *dictDb) GetCode(userId string) (*internal.Code, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db.codes[userId], nil
}

func (db *dictDb) GetNote(id string) (*internal.Note, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db.notes[id], nil
}

func (db *dictDb) GetToken(tokenId string) (*internal.Token, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	token, _ := db.tokens[tokenId]
	return token, nil
}

func (db *dictDb) GetUser(id string) (*internal.User, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db.users[id], nil
}

func (db *dictDb) IsAdmin(userId string) (bool, error) {
	return slices.Contains(db.admins, userId), nil
}

func (db *dictDb) IsNoteOwner(noteId string, principalId string) (bool, error) {
	n, ok := db.notes[noteId]

	if !ok {
		return false, DbError{"note does not exist, cannot be owner!!"}
	}

	return n.UserId == principalId, nil
}

func (db *dictDb) Initialize() error {
	db.admins = []string{"admin"}

	db.notes = map[string]*internal.Note{
		"dogs":    {Id: "dogs", Title: "Doggos", Text: "doggos are sweet", LastModified: 0, UserId: "tom"},
		"cats":    {Id: "cats", Title: "Cattos", Text: "meowowow", LastModified: 0, UserId: "tom"},
		"rabbits": {Id: "rabbits", Title: "Buns", Text: "hoppity hop, motherfuckas", LastModified: 0, UserId: "tom"},
		"goblins": {Id: "goblins", Title: "Gobbos", Text: "Grickle grackle", LastModified: 0, UserId: "mot"},
	}

	db.users = map[string]*internal.User{
		"tom":          {Id: "tom", Email: "fake@fake.com"},
		"mot":          {Id: "mot", Email: "ekaf@fake.com"},
		"expired-code": {Id: "expired-code", Email: "expired@fake.com"},
		"admin":        {Id: "admin", Email: "fakeadmin@fake.com"},
	}

	hash, err := HashString("expired")

	if err != nil {
		log.Fatal("Hash error during DB initialization!")
	}

	db.codes = map[string]*internal.Code{
		"expired-code": {"expired-code", hash, time.Unix(0, 0)},
	}

	db.tokens = map[string]*internal.Token{}

	db.initialized = true
	return nil
}

func (db *dictDb) ListLastModified(userId string) (map[string]int, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	data := make(map[string]int)

	for key := range maps.Keys(db.notes) {
		if db.notes[key].UserId == userId {
			data[key] = db.notes[key].LastModified
		}
	}

	return data, nil
}

func (db *dictDb) UpdateNote(newNote *internal.Note) error {
	if !db.initialized {
		return DbError{"Database not initialized!"}
	}

	oldNote, ok := db.notes[newNote.Id]

	if !ok {
		return DbError{"Could not find oldNote: " + newNote.Id}
	}

	oldNote.LastModified = newNote.LastModified
	oldNote.Text = newNote.Text
	oldNote.Title = newNote.Title
	return nil
}

func (db *dictDb) SetUser(u *internal.User) error {
	if !db.initialized {
		return DbError{"Database not initialized!"}
	}

	usr, ok := db.users[u.Id]

	if ok {
		usr.Email = u.Email
	} else {
		return DbError{"user not found in database!"}
	}

	return nil
}
