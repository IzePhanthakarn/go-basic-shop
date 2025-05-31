package entities

type Image struct {
	Id       string `db:"id" json:"id"`
	Filename string `db:"filename" json:"filename"`
	Url      string `db:"url" json:"url"`
}
