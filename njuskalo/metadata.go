package njuskalo

import (
	"encoding/json"
	"fmt"

	"github.com/sycomancy/glasnik/infra"
)

func headersForLocality() map[string]string {
	return map[string]string{
		"content-type":  "application/vnd.api+json",
		"user-agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"authorization": fmt.Sprintf("Bearer %s", infra.Config.NJUSKALO_BEARER_TOKEN),
	}
}

type LocalityEntry struct {
	Id         string `json:"id"`
	Attributes struct {
		Title string `json:"title"`
	}
}

type LocalityResponse struct {
	Data []LocalityEntry `json:"data"`
}

type NjuskaloMeta struct {
	client *infra.IncognitoClient
}

func NewNjuskaloMeta(client *infra.IncognitoClient) *NjuskaloMeta {
	return &NjuskaloMeta{
		client: client,
	}
}

func (m *NjuskaloMeta) RebuildLocalityMeta(locationIds []int) (*LocalityResponse, error) {
	data, err := m.fetchLocalityMeta(locationIds)
	if err != nil {
		return &LocalityResponse{}, err
	}

	entries := make([]interface{}, len(data.Data))
	for i := range data.Data {
		entries[i] = data.Data[i]
	}
	_ = infra.InsertDocuments("locality", entries)
	return data, nil
}

func (m *NjuskaloMeta) fetchLocalityMeta(locationIds []int) (*LocalityResponse, error) {
	url := "https://www.njuskalo.hr/ccapi/v2/locality?filter[parent]="
	for _, id := range locationIds {
		url += fmt.Sprint(id)
	}
	_, body, err := m.client.GetURLData(url, headersForLocality())
	if err != nil {
		return &LocalityResponse{}, err
	}

	var resp LocalityResponse
	err = json.NewDecoder(body).Decode(&resp)
	if err != nil {
		return &LocalityResponse{}, err
	}
	return &resp, nil
}
