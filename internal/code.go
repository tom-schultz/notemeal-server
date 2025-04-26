package internal

import (
	"encoding/json"
	"time"
)

const CodeJsonKey string = "code"

type Code struct {
	UserId     string
	CodeHash   string
	Expiration time.Time
}

func BuildCodeJson(code string) ([]byte, error) {
	inData := map[string]string{CodeJsonKey: code}
	outData, err := json.Marshal(inData)

	if err != nil {
		return nil, err
	}

	return outData, nil
}
