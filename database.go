package main

import (
	"maps"
)

type NotemealDb interface {
	DeleteNote(id string) error
	GetNote(id string) (*Note, error)
	Initialize() error
	ListLastModified() (map[string]int, error)
	SetNote(n *Note) error
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
	_users       map[string]*User
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

func (db *NotemealDictDb) Initialize() error {
	db._notes = map[string]*Note{
		"doggos":  {"doggos", "Doggos", "doggos are sweet", 0},
		"cattos":  {"cattos", "Cattos", "meowow", 0},
		"rabbits": {"rabbits", "Rabbits", "hoppity hop, mothafucka", 0},
	}

	db._initialized = true
	return nil
}

func (db *NotemealDictDb) ListLastModified() (map[string]int, error) {
	if !db._initialized {
		return nil, DbError{"Database not initialized!"}
	}

	data := make(map[string]int, len(db._notes))

	for key := range maps.Keys(db._notes) {
		data[key] = db._notes[key].LastModified
	}

	return data, nil
}

func (db *NotemealDictDb) SetNote(n *Note) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	db._notes[n.Id] = n
	return nil
}

func (db *NotemealDictDb) SetUser(u *User) error {
	if !db._initialized {
		return DbError{"Database not initialized!"}
	}

	db._users[u.Id] = u
	return nil
}
