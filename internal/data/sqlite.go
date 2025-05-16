package data

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"notemeal-server/internal"
	"os"
	"time"
)

type sqlite struct {
	db          *sql.DB
	initialized bool
}

func Sqlite(testing bool) Datasource {
	ds := &sqlite{}

	if err := ds.initialize(testing); err != nil {
		log.Fatal(err)
	}

	return ds
}

func (ds *sqlite) DeleteCode(userId string) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`delete from code where id=?`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(userId)
	return err
}

func (ds *sqlite) DeleteNote(id string) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`delete from note where id=?`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

func (ds *sqlite) DeleteUser(id string) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`delete from user where id=?`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(id)
	return err
}

func (ds *sqlite) GetCode(userId string) (*internal.Code, error) {
	if !ds.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`select id, hash, expiration from code where id=?`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var code internal.Code
	err = stmt.QueryRow(userId).Scan(&code.UserId, &code.Hash, &code.Expiration)

	if err == nil {
		return &code, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (ds *sqlite) GetNote(id string) (*internal.Note, error) {
	if !ds.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`select * from note where id=?`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var note internal.Note
	err = stmt.QueryRow(id).Scan(&note.Id, &note.LastModified, &note.Text, &note.Title, &note.UserId)

	if err == nil {
		return &note, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (ds *sqlite) GetNotesByUser(userId string) ([]*internal.Note, error) {
	if !ds.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	data := make([]*internal.Note, 0)

	stmt, err := ds.db.Prepare(`select * from note where userId=?`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	rows, err := stmt.Query(userId)
	var note *internal.Note

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		note = &internal.Note{}
		err = rows.Scan(&note.Id, &note.LastModified, &note.Text, &note.Title, &note.UserId)

		if err != nil {
			return nil, err
		}

		data = append(data, note)
	}

	return data, nil
}

func (ds *sqlite) userHasToken(userId string) (bool, error) {
	stmt, err := ds.db.Prepare(`select id, hash, userId from token where userId=?`)

	if err != nil {
		return false, err
	}

	defer stmt.Close()
	row, err := stmt.Query(userId)

	if err != nil {
		return false, err
	}

	return row.Next(), nil
}

func (ds *sqlite) GetToken(id string) (*internal.Token, error) {
	if !ds.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}
	stmt, err := ds.db.Prepare(`select id, hash, userId from token where id=?`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var token internal.Token
	row := stmt.QueryRow(id)
	err = row.Scan(&token.Id, &token.Hash, &token.UserId)

	if err == nil {
		return &token, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (ds *sqlite) GetUser(id string) (*internal.User, error) {
	if !ds.initialized {
		return nil, internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`select * from user where id=?`)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()
	var user internal.User
	err = stmt.QueryRow(id).Scan(&user.Id, &user.Email, &user.IsAdmin)

	if err == nil {
		return &user, nil
	} else if err == sql.ErrNoRows {
		return nil, nil
	} else {
		return nil, err
	}
}

func (ds *sqlite) initialize(testing bool) error {
	dbPath := "./notemeal.db"

	if testing {
		dbPath = "./notemeal-test.db"
		err := os.Remove(dbPath)

		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Fatal(err)
		}
	}

	db, err := sql.Open("sqlite3", dbPath)

	if err != nil {
		return err
	}

	ds.db = db

	err = ds.initializeTokens(testing)

	if err != nil {
		return err
	}

	err = ds.initializeCodes(testing)

	if err != nil {
		return err
	}

	err = ds.initializeNotes(testing)

	if err != nil {
		return err
	}

	err = ds.initializeUsers(testing)

	if err != nil {
		return err
	}

	ds.initialized = true
	return nil
}

func (ds *sqlite) initializeCodes(testing bool) error {
	sqlStmt := `
	create table if not exists code (id text not null primary key, hash text, expiration datetime);
	`

	_, err := ds.db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	codes := map[string]*internal.Code{}

	if !testing {
		adminTokenExists, err := ds.userHasToken("admin")

		if err != nil {
			return err
		}

		if adminTokenExists {
			return nil
		}

		adminCode := internal.CreateSecureString()
		adminHash, err := internal.HashString(adminCode)
		fmt.Printf("Admin code: %s\n", adminCode)

		if err != nil {
			return err
		}

		codes = map[string]*internal.Code{
			"admin": {"admin", adminHash, time.Now().Add(24 * time.Hour)},
		}
	} else {
		expiredHash, err := internal.HashString("expired")

		if err != nil {
			return nil
		}

		turtleHash, err := internal.HashString("turtles")

		if err != nil {
			return nil
		}

		codes = map[string]*internal.Code{
			"expired-code": {"expired-code", expiredHash, time.Unix(0, 0)},
			"tom":          {"tom", turtleHash, time.Now().Add(time.Hour)},
		}
	}

	tx, err := ds.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("insert or replace into code(id, hash, expiration) values(?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	for _, code := range codes {
		_, err = stmt.Exec(code.UserId, code.Hash, code.Expiration)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (ds *sqlite) initializeNotes(testing bool) error {
	sqlStmt := `
	create table if not exists note (id text not null primary key, lastModified integer not null, text text not null, title text not null, userId text not null);
	`

	_, err := ds.db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	if !testing {
		return nil
	}

	tx, err := ds.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`insert or replace into note(id, lastModified, text, title, userId) values(?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	notes := map[string]*internal.Note{
		"dogs":    {Id: "dogs", Title: "Doggos", Text: "doggos are sweet", LastModified: 0, UserId: "tom"},
		"cats":    {Id: "cats", Title: "Cattos", Text: "meowowow", LastModified: 0, UserId: "tom"},
		"rabbits": {Id: "rabbits", Title: "Buns", Text: "hoppity hop, motherfuckas", LastModified: 0, UserId: "tom"},
		"goblins": {Id: "goblins", Title: "Gobbos", Text: "Grickle grackle", LastModified: 0, UserId: "mot"},
	}

	for _, note := range notes {
		_, err = stmt.Exec(note.Id, note.LastModified, note.Text, note.Title, note.UserId)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (ds *sqlite) initializeUsers(testing bool) error {
	sqlStmt := `
	create table if not exists user (id text not null primary key, email text not null, isAdmin integer not null);
	`

	_, err := ds.db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return err
	}

	if !testing {
		return nil
	}

	tx, err := ds.db.Begin()

	if err != nil {
		return err
	}

	stmt, err := tx.Prepare(`insert or replace into user(id, email, isAdmin) values(?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()

	users := map[string]*internal.User{
		"tom":          {Id: "tom", Email: "fake@fake.com"},
		"mot":          {Id: "mot", Email: "ekaf@fake.com"},
		"expired-code": {Id: "expired-code", Email: "expired@fake.com"},
		"admin":        {Id: "admin", Email: "fakeadmin@fake.com", IsAdmin: true},
	}

	for _, user := range users {
		_, err = stmt.Exec(user.Id, user.Email, user.IsAdmin)

		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	return err
}

func (ds *sqlite) initializeTokens(testing bool) error {
	sqlStmt := `
	create table if not exists token (id text not null primary key, hash text not null, userId text not null);
	`

	_, err := ds.db.Exec(sqlStmt)

	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
	}

	if !testing {
		return nil
	}

	return err
}

func (ds *sqlite) UpdateCode(code *internal.Code) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`insert or replace into code(id, hash, expiration) values(?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(code.UserId, code.Hash, code.Expiration)
	return err
}

func (ds *sqlite) UpdateNote(note *internal.Note) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`insert or replace into note(id, lastModified, text, title, userId) values(?, ?, ?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(note.Id, note.LastModified, note.Text, note.Title, note.UserId)
	return err
}

func (ds *sqlite) UpdateToken(token *internal.Token) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`insert or replace into token(id, hash, userId) values(?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(token.Id, token.Hash, token.UserId)
	return err
}

func (ds *sqlite) UpdateUser(user *internal.User) error {
	if !ds.initialized {
		return internal.Error{"Datasource not initialized!"}
	}

	stmt, err := ds.db.Prepare(`insert or replace into user(id, email, isAdmin) values(?, ?, ?)`)

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(user.Id, user.Email, user.IsAdmin)
	return err
}
