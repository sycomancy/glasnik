package types

type AddPageResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
	Price int32  `json:"price"`
}
