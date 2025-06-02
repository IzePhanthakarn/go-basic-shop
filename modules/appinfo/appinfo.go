package appinfo

type CategoryFilter struct {
	Title string `query:"title"`
}

type Category struct {
	Id    int    `db:"id" json:"id"`
	Title string `db:"title" json:"title"`
}

type GenerateApiKeyRes struct {
	ApiKey string `json:"api_key"`
}
type CategoryRemoveRes struct {
	CategoryId int `json:"category_id"`
}