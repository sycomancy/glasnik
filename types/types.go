package types

type AdsPageResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
	Price string `json:"price"`
}

type RequestData struct {
	Filter string `json:"filter"`
}
