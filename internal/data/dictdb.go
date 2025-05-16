package data

import (
	"log"
	"maps"
	"notemeal-server/internal"
	"time"
)

type dictDb struct {
	admins      []string
	notes       map[string]*internal.Note
	tokens      map[string]*internal.Token
	codes       map[string]*internal.Code
	users       map[string]*internal.User
	initialized bool
}

func DictDb() Datasource {
	db := &dictDb{}

	if err := db.initialize(); err != nil {
		log.Fatal(err)
	}

	return db
}

func (db *dictDb) DeleteCode(userId string) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	delete(db.codes, userId)
	return nil
}

func (db *dictDb) DeleteNote(id string) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	delete(db.notes, id)
	return nil
}

func (db *dictDb) DeleteUser(id string) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	delete(db.users, id)
	return nil
}

func (db *dictDb) GetCode(userId string) (*internal.Code, error) {
	if !db.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	return db.codes[userId], nil
}

func (db *dictDb) GetNote(id string) (*internal.Note, error) {
	if !db.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	return db.notes[id], nil
}

func (db *dictDb) GetNotesByUser(userId string) ([]*internal.Note, error) {
	if !db.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	data := make([]*internal.Note, 0)

	for key := range maps.Keys(db.notes) {
		if db.notes[key].UserId == userId {
			data = append(data, db.notes[key])
		}
	}

	return data, nil
}

func (db *dictDb) GetToken(id string) (*internal.Token, error) {
	if !db.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	return db.tokens[id], nil
}

func (db *dictDb) GetUser(id string) (*internal.User, error) {
	if !db.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	return db.users[id], nil
}

func (db *dictDb) initialize() error {
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
		"admin":        {Id: "admin", Email: "fakeadmin@fake.com", IsAdmin: true},
	}
	expiredHash, err := internal.HashString("expired")

	if err != nil {
		return nil
	}

	codeHash, err := internal.HashString("turtles")

	if err != nil {
		return nil
	}

	db.codes = map[string]*internal.Code{
		"expired-code": {"expired-code", expiredHash, time.Unix(0, 0)},
		"tom":          {"tom", codeHash, time.Now().Add(time.Hour)},
	}

	db.tokens = map[string]*internal.Token{}
	db.initialized = true
	return nil
}

func (db *dictDb) UpdateCode(code *internal.Code) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	db.codes[code.UserId] = code
	return nil
}

func (db *dictDb) UpdateNote(note *internal.Note) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	db.notes[note.Id] = note
	return nil
}

func (db *dictDb) UpdateToken(token *internal.Token) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	db.tokens[token.Id] = token
	return nil
}

func (db *dictDb) UpdateUser(user *internal.User) error {
	if !db.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	db.users[user.Id] = user
	return nil
}
