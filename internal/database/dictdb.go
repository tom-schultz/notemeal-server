package database

import (
	"log"
	"maps"
	"notemeal-server/internal"
	"time"
)

type dictDb struct {
	_initialized bool
	_notes       map[string]*internal.Note
	_tokens      map[string]*internal.Token
	_Codes       map[string]*internal.Code
	_users       map[string]*internal.User
}

func DictDb() {
	Db = &dictDb{}

	if err := Db.Initialize(); err != nil {
		log.Fatal(err)
	}
}

func (db *dictDb) CreateOrUpdateCode(id string) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	expiration := time.Now().Add(time.Hour)
	code := "blue cat dog"
	codeHash := HashString(code)

	if code, ok := db._Codes[id]; ok {
		code.Expiration = expiration
	} else {
		db._Codes[id] = &internal.Code{
			UserId:     id,
			CodeHash:   codeHash,
			Expiration: expiration,
		}
	}

	return nil
}

func (db *dictDb) CreateToken(userId string, CodeString string) (string, error) {
	if !db._initialized {
		return "", DbError{"Database not initialized!"}
	}

	hashedCodeString := HashString(CodeString)
	Code := db._Codes[userId]

	if !CompareHashedString(hashedCodeString, Code.CodeHash) || Code.Expiration.Before(time.Now()) {
		return "", DbError{"Invalid tkn code!"}
	}

	tkn := "123456"
	tokenHash := HashString(tkn)

	db._tokens[tokenHash] = &internal.Token{TokenHash: tokenHash, UserId: userId}

	return tkn, nil
}

func CompareHashedString(str string, hashedStr string) bool {
	return str == hashedStr
}

func HashString(str string) string {
	return "!!" + str
}

func (db *dictDb) DeleteUser(id string) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	delete(db._users, id)
	return nil
}

func (db *dictDb) DeleteNote(id string) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	delete(db._notes, id)
	return nil
}

func (db *dictDb) GetNote(id string) (*internal.Note, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db._notes[id], nil
}

func (db *dictDb) GetToken(token string) (*internal.Token, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	tokenHash := HashString(token)
	return db._tokens[tokenHash], nil
}

func (db *dictDb) GetUser(id string) (*internal.User, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	return db._users[id], nil
}

func (db *dictDb) IsNoteOwner(noteId string, principalId string) (bool, error) {
	n, ok := db._notes[noteId]

	if !ok {
		return false, DbError{"note does not exist, cannot be owner!!"}
	}

	return n.UserId == principalId, nil
}

func (db *dictDb) Initialize() error {
	db._notes = map[string]*internal.Note{
		"dogs":    {Id: "dogs", Title: "Doggos", Text: "doggos are sweet", LastModified: 0, UserId: "tom"},
		"cats":    {Id: "cats", Title: "Cattos", Text: "meowowow", LastModified: 0, UserId: "tom"},
		"rabbits": {Id: "rabbits", Title: "Buns", Text: "hoppity hop, motherfuckas", LastModified: 0, UserId: "tom"},
		"goblins": {Id: "goblins", Title: "Gobbos", Text: "Grickle grackle", LastModified: 0, UserId: "mot"},
	}

	db._users = map[string]*internal.User{
		"tom":   {"tom", "fake@fake.com"},
		"mot":   {"mot", "ekaf@fake.com"},
		"admin": {"admin", "fakeadmin@fake.com"},
	}

	db._Codes = map[string]*internal.Code{}
	db._tokens = map[string]*internal.Token{}

	db._initialized = true
	return nil
}

func (db *dictDb) ListLastModified(userId string) (map[string]int, error) {
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

func (db *dictDb) UpdateNote(newNote *internal.Note) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	oldNote, ok := db._notes[newNote.Id]

	if !ok {
		return DbError{"Could not find oldNote: " + newNote.Id}
	}

	oldNote.LastModified = newNote.LastModified
	oldNote.Text = newNote.Text
	oldNote.Title = newNote.Title
	return nil
}

func (db *dictDb) SetUser(u *internal.User) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	usr, ok := db._users[u.Id]

	if ok {
		usr.Email = u.Email
	} else {
		return DbError{"user not found in database!"}
	}

	return nil
}
