package njuskalo

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

const entityListItemsQuery = ".EntityList--Standard > .EntityList-items > .EntityList-item.EntityList-item--Regular"
const entityTitleQuery = ".entity-title a"
const entityPriceQuery = ".price-item > .price"

var headers = map[string]string{"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"}

var ErrBadRequest = fmt.Errorf("400 Bad Request")

func Fetch(url string, client *infra.IncognitoClient) ([]types.AdsPageResponse, error) {
	items := make([]types.AdsPageResponse, 0)

	hasMorePage := true
	page := 1

	for hasMorePage {
		pageItems, err := getItemsForPage(url, page, client)

		if err != nil {
			return nil, err
		}

		if len(pageItems) == 0 {
			hasMorePage = false
		} else {
			items = append(items, pageItems...)
			hasMorePage = true
			page += 1
		}
	}

	if len(items) != 0 {
		page += 1
	}

	return items, nil
}

func getItemsForPage(url string, page int, client *infra.IncognitoClient) ([]types.AdsPageResponse, error) {
	urlWithPage := url

	if page != 1 {
		urlWithPage = fmt.Sprintf("%s%s%d", url, "&page=", page)
	}

	status, body, error := client.GetURLDataWithRetries(urlWithPage, headers)
	// Find review items
	items := []types.AdsPageResponse{}

	if status == ErrBadRequest.Error() {
		// TODO(sycomancy): handle 400 BAD REQUEST!!!
		logrus.WithFields(logrus.Fields{
			"status": status,
		}).Warning("njuskalo returned 400 bad request")
		return nil, ErrBadRequest
	}

	if strings.Contains(status, "404") {
		return items, nil
	}

	if error != nil {
		return nil, error
	}

	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return nil, fmt.Errorf("unable to get document from response %w", err)
	}

	doc.Find(entityListItemsQuery).Each(func(i int, s *goquery.Selection) {
		var id string
		var title string
		var link string
		var price string

		itemId, idExists := s.Attr("data-href")
		titleEl := s.Find(entityTitleQuery)
		hrefAttr, exists := titleEl.Attr("href")
		price = s.Find(entityPriceQuery).Text()
		title = titleEl.Text()

		if idExists {
			id = itemId
		}

		if exists {
			link = hrefAttr
		}

		item := types.AdsPageResponse{Id: id, Title: title, Link: link, Price: price}
		items = append(items, item)
	})

	return items, nil
}
