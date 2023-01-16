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
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	locationMetaCollection = "locality"
	baseLocalityUrl        = "https://www.njuskalo.hr/ccapi/v2/locality?filter[parent]="
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "locations",
})

type LocalityEntry struct {
	Id    primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Title string             `json:"title"`
}

type LocalityResponse struct {
	Data []LocalityEntry `json:"data"`
}

type LocationPageResult struct {
	Location  LocalityEntry   `bson:"location,omitempty"`
	Page      int             `bson:"page,omitempty"`
	Items     []types.AdEntry `bson:"items,omitempty"`
	Err       error           `bson:"err,omitempty"`
	Completed bool            `bson:"completed,omitempty"`
}

func GetAllLocalityEntries() []*LocalityEntry {
	var locations []*LocalityEntry
	infra.FindDocuments("locality", bson.D{}, &locations)
	return locations
}

type LocationService struct {
}

func NewLocationService() *LocationService {
	return &LocationService{}
}

func (l *LocationService) GetLocationPages(loc *LocalityEntry, result chan *LocationPageResult, client *infra.IncognitoClient) {
	page := 1
	hasMorePage := true
	locationURL := fmt.Sprintf("%s%s", baseURL, loc.Id)

	for hasMorePage {
		locationPageHTML, err := l.getPageHTML(locationURL, page, client)
		if err != nil {
			result <- &LocationPageResult{
				Location:  *loc,
				Err:       err,
				Page:      page,
				Completed: false,
			}
			hasMorePage = false
			break
		}

		adsForPage, err := l.getItemsFromHTML(locationPageHTML)
		if err != nil {
			flogg.Error("Unable to parse location page html for %s page: %d", loc.Title, page)
			result <- &LocationPageResult{
				Location:  *loc,
				Err:       err,
				Page:      page,
				Completed: false,
			}
			hasMorePage = false
			break
		}

		hasMorePage = len(adsForPage) > 0

		result <- &LocationPageResult{
			Location:  *loc,
			Page:      page,
			Items:     adsForPage,
			Completed: !hasMorePage,
		}

		page += 1
	}
}

func (l *LocationService) getPageHTML(url string, page int, client *infra.IncognitoClient) (string, error) {
	// TODO(sycomancy): inversion of control
	pageURL, err := njuskalo.GetUrlForPage(url, page)
	if err != nil {
		return "", err
	}

	_, res, err := client.GetURLDataWithRetries(pageURL, njuskalo.Headers)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	buf.ReadFrom(res)
	return buf.String(), nil
}

func (l *LocationService) getItemsFromHTML(html string) ([]types.AdEntry, error) {
	return njuskalo.ParsePage(html)
}

func (l *LocationService) FetchAndPersistLocationMeta(locationIds []string, client *infra.IncognitoClient) error {
	data, err := l.fetchLocalityMeta(locationIds, client)
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

func (l *LocationService) fetchLocalityMeta(locationIds []string, client *infra.IncognitoClient) (*LocalityResponse, error) {
	url := baseLocalityUrl
	for _, id := range locationIds {
		url += fmt.Sprintf("%s,", id)
	}
	_, body, err := client.GetURLData(url, l.generateHeadersForMeta())
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

func (l *LocationService) generateHeadersForMeta() map[string]string {
	return map[string]string{
		"content-type":  "application/vnd.api+json",
		"user-agent":    "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36",
		"authorization": fmt.Sprintf("Bearer %s", infra.Config.NJUSKALO_BEARER_TOKEN),
	}
}
