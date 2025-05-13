package internal

type Note struct {
	Id           string `json:"id"`
	LastModified int    `json:"lastModified"`
	Text         string `json:"text"`
	Title        string `json:"title"`
	UserId       string `json:"-"`
}
