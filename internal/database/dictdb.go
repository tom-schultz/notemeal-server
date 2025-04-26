package database

import (
	"log"
	"maps"
	"notemeal-server/internal"
	"time"
)

type dictDb struct {
	initialized bool
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

func (db *dictDb) CreateOrUpdateCode(id string) (string, error) {
	if !db.initialized {
		return "", DbError{"Database not initialized!"}
	}

	expiration := time.Now().Add(time.Hour)
	code := "blue cat dog"
	codeHash := HashString(code)

	if code, ok := db.codes[id]; ok {
		code.Expiration = expiration
	} else {
		db.codes[id] = &internal.Code{
			UserId:     id,
			CodeHash:   codeHash,
			Expiration: expiration,
		}
	}

	return code, nil
}

func (db *dictDb) CreateToken(userId string, CodeString string) (string, error) {
	if !db.initialized {
		return "", DbError{"Database not initialized!"}
	}

	hashedCodeString := HashString(CodeString)
	code := db.codes[userId]

	if code == nil {
		return "", nil
	}

	if !CompareHashedString(hashedCodeString, code.CodeHash) || code.Expiration.Before(time.Now()) {
		return "", nil
	}

	token := "123456"
	tokenHash := HashString(token)

	db.tokens[tokenHash] = &internal.Token{TokenHash: tokenHash, UserId: userId}
	delete(db.codes, userId)

	return token, nil
}

func CompareHashedString(str string, hashedStr string) bool {
	return str == hashedStr
}

func HashString(str string) string {
	return "!!" + str
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

func (db *dictDb) GetNote(id string) (*internal.Note, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db.notes[id], nil
}

func (db *dictDb) GetToken(token string) (*internal.Token, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	tokenHash := HashString(token)
	return db.tokens[tokenHash], nil
}

func (db *dictDb) GetUser(id string) (*internal.User, error) {
	if !db.initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db.users[id], nil
}

func (db *dictDb) IsNoteOwner(noteId string, principalId string) (bool, error) {
	n, ok := db.notes[noteId]

	if !ok {
		return false, DbError{"note does not exist, cannot be owner!!"}
	}

	return n.UserId == principalId, nil
}

func (db *dictDb) Initialize() error {
	db.notes = map[string]*internal.Note{
		"dogs":    {Id: "dogs", Title: "Doggos", Text: "doggos are sweet", LastModified: 0, UserId: "tom"},
		"cats":    {Id: "cats", Title: "Cattos", Text: "meowowow", LastModified: 0, UserId: "tom"},
		"rabbits": {Id: "rabbits", Title: "Buns", Text: "hoppity hop, motherfuckas", LastModified: 0, UserId: "tom"},
		"goblins": {Id: "goblins", Title: "Gobbos", Text: "Grickle grackle", LastModified: 0, UserId: "mot"},
	}

	db.users = map[string]*internal.User{
		"tom":   {Id: "tom", Email: "fake@fake.com"},
		"mot":   {Id: "mot", Email: "ekaf@fake.com"},
		"admin": {Id: "admin", Email: "fakeadmin@fake.com"},
	}

	db.codes = map[string]*internal.Code{}
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
