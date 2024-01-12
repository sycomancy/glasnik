package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type AdInDb struct {
	ID   primitive.ObjectID `bson:"_id"`
	Slug string             `bson:"slug"`
	Add  AdEntry
}

// Actual result from web scrapping
type AdEntry struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Link        string `json:"link"`
	Price       string `json:"price"`
	Description string `json:"description"`
}

type BasicInfoEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type AdEntryDetails struct {
	Description string `json:"description"`
	Attributes  []BasicInfoEntry
}

type FetchEntryDetailsResult struct {
	Ad      AdInDb
	Details AdEntryDetails
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
