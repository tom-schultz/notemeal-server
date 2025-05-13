package internal

const TokenJsonKey string = "token"

type Token struct {
	Id     string
	UserId string
	Hash   string
}

type ClientToken struct {
	Id    string `json:"id"`
	Token string `json:"token"`
}
