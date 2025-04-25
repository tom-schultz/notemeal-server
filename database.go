package main

import (
	"maps"
	"time"
)

type NotemealDb interface {
	CreateOrUpdateTokenCode(id string) error
	CreateToken(userId string, tokenCodeString string) (string, error)
	DeleteNote(id string) error
	DeleteUser(id string) error
	GetNote(id string) (*Note, error)
	GetToken(token string) (*Token, error)
	GetUser(id string) (*User, error)
	IsNoteOwner(noteId string, principalId string) (bool, error)
	Initialize() error
	ListLastModified(userId string) (map[string]int, error)
	UpdateNote(n *Note) error
	SetUser(u *User) error
}

type DbError struct {
	msg string
}

func (e DbError) Error() string {
	return e.msg
}

type NotemealDictDb struct {
	_initialized bool
	_notes       map[string]*Note
	_tokens      map[string]*Token
	_tokenCodes  map[string]*TokenCode
	_users       map[string]*User
}

func (db *NotemealDictDb) CreateOrUpdateTokenCode(id string) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	expiration := time.Now().Add(time.Hour)
	code := "blue cat dog"
	codeHash := HashString(code)

	if code, ok := db._tokenCodes[id]; ok {
		code.Expiration = expiration
	} else {
		db._tokenCodes[id] = &TokenCode{
			UserId:     id,
			CodeHash:   codeHash,
			Expiration: expiration,
		}
	}

	return nil
}

func (db *NotemealDictDb) CreateToken(userId string, tokenCodeString string) (string, error) {
	if !db._initialized {
		return "", DbError{"Database not initialized!"}
	}

	hashedTokenCodeString := HashString(tokenCodeString)
	tokenCode := db._tokenCodes[userId]

	if !CompareHashedString(hashedTokenCodeString, tokenCode.CodeHash) || tokenCode.Expiration.Before(time.Now()) {
		return "", DbError{"Invalid token code!"}
	}

	token := "123456"
	tokenHash := HashString(token)

	db._tokens[tokenHash] = &Token{TokenHash: tokenHash, UserId: userId}

	return token, nil
}

func CompareHashedString(str string, hashedStr string) bool {
	return str == hashedStr
}

func HashString(str string) string {
	return "!!" + str
}

func (db *NotemealDictDb) DeleteUser(id string) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	delete(db._users, id)
	return nil
}

func (db *NotemealDictDb) DeleteNote(id string) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	delete(db._notes, id)
	return nil
}

func (db *NotemealDictDb) GetNote(id string) (*Note, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db._notes[id], nil
}

func (db *NotemealDictDb) GetToken(token string) (*Token, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	tokenHash := HashString(token)
	return db._tokens[tokenHash], nil
}

func (db *NotemealDictDb) GetUser(id string) (*User, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db._users[id], nil
}

func (db *NotemealDictDb) IsNoteOwner(noteId string, principalId string) (bool, error) {
	note, ok := db._notes[noteId]

	if !ok {
		return false, DbError{"Note does not exist, cannot be owner!!"}
	}

	return note.UserId == principalId, nil
}

func (db *NotemealDictDb) Initialize() error {
	db._notes = map[string]*Note{
		"dogs":    {Id: "dogs", Title: "Doggos", Text: "doggos are sweet", LastModified: 0, UserId: "tom"},
		"cats":    {Id: "cats", Title: "Cattos", Text: "meowowow", LastModified: 0, UserId: "tom"},
		"rabbits": {Id: "rabbits", Title: "Buns", Text: "hoppity hop, motherfuckas", LastModified: 0, UserId: "tom"},
		"goblins": {Id: "goblins", Title: "Gobbos", Text: "Grickle grackle", LastModified: 0, UserId: "mot"},
	}

	db._users = map[string]*User{
		"tom":   {"tom", "fake@fake.com"},
		"mot":   {"mot", "ekaf@fake.com"},
		"admin": {"admin", "fakeadmin@fake.com"},
	}

	db._tokenCodes = map[string]*TokenCode{}
	db._tokens = map[string]*Token{}

	db._initialized = true
	return nil
}

func (db *NotemealDictDb) ListLastModified(userId string) (map[string]int, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	data := make(map[string]int)

	for key := range maps.Keys(db._notes) {
		if db._notes[key].UserId == userId {
			data[key] = db._notes[key].LastModified
		}
	}

	return data, nil
}

func (db *NotemealDictDb) UpdateNote(newNote *Note) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	note, ok := db._notes[newNote.Id]

	if !ok {
		return DbError{"Could not find note: " + newNote.Id}
	}

	note.LastModified = newNote.LastModified
	note.Text = newNote.Text
	note.Title = newNote.Title
	return nil
}

func (db *NotemealDictDb) SetUser(u *User) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	user, ok := db._users[u.Id]

	if ok {
		user.Email = u.Email
	} else {
		return DbError{"User not found in database!"}
	}

	return nil
}
