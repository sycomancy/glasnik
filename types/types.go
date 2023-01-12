package types

type AdDetails struct {
	Size  int `json:"size" bson:"size"`
	Rooms int `json:"rooms" bson:"rooms"`
}

// Actual result from web scrapping
type AdEntry struct {
	Id    string `json:"id" bson:"id"`
	Title string `json:"title" bson:"title,omitempty"`
	Link  string `json:"link" bson:"link,omitempty"`
	Price string `json:"price" bson:"price,omitempty"`
	Size  string `json:"size" bson:"size,omitempty"`
}

// Request for data fetching.
// If Callback is provided, result will be delivered async with PUT on callback URL
type Request struct {
	Filter      string `json:"filter"`
	Token       string `json:"token"`
	CallbackURL string `json:"callbackUrl"`
}

// Result for data fetching request
// Data should be available is request is sync
type RequestResult struct {
	Data        []AdEntry `json:"data,omitempty"`
	Status      string    `json:"status"`
	CallbackURL string    `json:"callbackUrl,omitempty"`
	RequestID   int       `json:"requestId"`
}

// Service input param, maybe remove from public types
type RequestData struct {
	Filter      string
	Token       string
	CallbackURL string
	RequestID   int
}
