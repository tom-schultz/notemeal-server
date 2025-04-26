package database

type DbError struct {
	msg string
}

func (e DbError) Error() string {
	return e.msg
}
