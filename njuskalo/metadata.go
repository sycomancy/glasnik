package njuskalo

import (
	"encoding/json"
	"fmt"

	"github.com/sycomancy/glasnik/infra"
)

// TODO(sycomancy): change this
var BearerToken string = "Bearer eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJhdWQiOiJuanVza2Fsb19qc19hcHAiLCJqdGkiOiJlODg4ZjIwODllMTdhY2ZmMjEyMzc2YWU1M2U4Yzg5ZjUzZTBiOWIzMDhkOWEwM2U1YjFlYmMwOGNiZTFkMmFkOTAxOWQ4YTU0MDMzNjI0MCIsImlhdCI6MTY3MTEzMDM3MSwibmJmIjoxNjcxMTMwMzcxLCJleHAiOjE2NzExNTE5NzEsInN1YiI6IiIsInNjb3BlcyI6W119.elcUlwKPT0qvMN6M4eCfKsXPZRaJF2MvYCMqsntacD3yqOXwrNaC240Vuu6Xje3jCzCPJ16D-Na8zOouFgkU7AgKgd4azd7feVJSEq7Dc5Bjmz_14qBMYH5SZISxsMFdL00ZbYmgE0I2v7vTFCeeLIQ5nABhZlvdfwrPAOyQ67Z6Zf0t3rA0W9jYx0LKtlVSPMnZX0NkHxz3xLoy3hooABqLPn3GgwcBtytTToP6UJef_EQgtMn3eFvNPFpv6QESUdkA1cYmE9cbu8sx2XIc1REtNoHItZNRfH2HNQ0DcAnHqu7F8rV1guUd2Y3EoWwjzU_0pcO51vBHPaxPSAI1pQ"

var headersForLocality = map[string]string{
	"user-agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
	"authorization": BearerToken,
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
	_, body, err := m.client.GetURLData(url, headersForLocality)
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
