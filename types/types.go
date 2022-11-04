package types

// Actual result from web scrapping
type AdsPageResponse struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Link  string `json:"link"`
	Price string `json:"price"`
}

// Request for data fetching.
// If Callback is provided, result will be delivered async with PUT on callback URL
type RequestDataDTO struct {
	Filter      string `json:"filter"`
	Token       string `json:"token"`
	CallbackURL string `json:"callback"`
}

// Result for data fetching request
// Data should be available is request is sync
type RequestResult struct {
	Data        []AdsPageResponse
	Status      string
	CallbackURL string
}

// Service input param, maybe remove from public types
type RequestData struct {
	RequestDataDTO
	RequestID string
}
