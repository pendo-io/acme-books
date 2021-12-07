package models

type Book struct {
	Id       int64
	Title    string `json:"title"`
	Author   string `json:"author"`
	Borrowed bool   `json:"borrowed"`
}
