package service

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
	"github.com/sycomancy/glasnik/types"
	"go.mongodb.org/mongo-driver/bson"
)

var (
	locationMetaCollection = "locality"
	baseLocalityUrl        = "https://www.njuskalo.hr/ccapi/v2/locality?filter[parent]="
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "locations",
})

type LocalityEntry struct {
	Id         string `json:"id,omitempty" bson:"_id"`
	Attributes struct {
		Title string `json:"title"`
	}
}

type LocalityResponse struct {
	Data []LocalityEntry `json:"data"`
}

func GetAllLocalityEntries() []*LocalityEntry {
	var locations []*LocalityEntry
	infra.FindDocuments("locality", bson.D{}, &locations)
	return locations
}

type Location struct {
	client *infra.IncognitoClient
}

func NewLocation(client *infra.IncognitoClient) *Location {
	return &Location{client: client}
}

func (l *Location) GetPageHTML(url string, page int, client *infra.IncognitoClient) (string, error) {
	// TODO(sycomancy): inversion of control
	pageURL, err := njuskalo.GetUrlForPage(url, page)
	if err != nil {
		return "", err
	}

	_, res, err := client.GetURLDataWithRetries(pageURL, njuskalo.Headers)
	//errors.Is(err,njuskalo.ErrBadRequest)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	return buf.String(), nil
}

func (l *Location) GetItemsFromHTML(html string) ([]types.AdEntry, error) {
	return njuskalo.ParsePage(html)
}

// Location MetaData
func (l *Location) FetchAndPersistLocationMeta(locationIds []string) error {
	data, err := l.fetchLocalityMeta(locationIds)
	if err != nil {
		return err
	}
	entries := make([]interface{}, len(data.Data))
	for i := range data.Data {
		entries[i] = data.Data[i]
	}
	_ = infra.InsertDocuments(locationMetaCollection, entries)
	logg.Infof("Fetched and persisted data for %d locations /n", len(entries))
	return nil
}

func (l *Location) fetchLocalityMeta(locationIds []string) (*LocalityResponse, error) {
	url := baseLocalityUrl
	for _, id := range locationIds {
		url += fmt.Sprintf("%s,", id)
	}
	_, body, err := l.client.GetURLData(url, l.generateHeadersForMeta())
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

func (l *Location) generateHeadersForMeta() map[string]string {
	return map[string]string{
		"content-type":  "application/vnd.api+json",
		"user-agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"authorization": fmt.Sprintf("Bearer %s", infra.Config.NJUSKALO_BEARER_TOKEN),
	}
}
