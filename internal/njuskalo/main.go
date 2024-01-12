package njuskalo

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/types"
)

// For listing
const entityListItemsQuery = ".EntityList--Standard > .EntityList-items > .EntityList-item.EntityList-item--Regular"
const entityTitleQuery = ".entity-title a"
const entityPriceQuery = ".price-item > .price"
const descriptionQuery = ".entity-description-main"

// for details
const entityDetailsQuery = ".BlockStandard.ClassifiedDetailDescription > div > div"

// body > div.wrap-main > div.wrap-content.ClassifiedDetail.cf > main > article > div.content-primary > div > div.content-main > div.BlockStandard.ClassifiedDetailDescription > div > div

var headers = map[string]string{"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"}

var ErrBadRequest = fmt.Errorf("400 Bad Request")

func FetchEntryDetails(entries []types.AdInDb, result chan<- types.FetchEntryDetailsResult, client *infra.IncognitoClient) {
	for _, entry := range entries {
		itemDetails, err := getDetailsForPage(entry.Add.Link, client)
		if err != nil {
			close(result)
			return
		}
		result <- types.FetchEntryDetailsResult{Ad: entry, Details: itemDetails}
	}
}

func FetchEntries(url string, result chan<- []types.AdEntry, client *infra.IncognitoClient) {
	hasMorePage := true
	page := 1
	items := make([]types.AdEntry, 0)

	for hasMorePage {
		pageItems, err := getItemsForPage(url, page, client)

		if err != nil {
			fmt.Println(err)
			close(result)
			return
		}

		if len(pageItems) == 0 {
			hasMorePage = false
		} else {
			items = pageItems
			hasMorePage = true
			page += 1
			result <- items
		}
	}

	if len(items) != 0 {
		page += 1
	}

	close(result)
	return
}

func Fetch(url string, client *infra.IncognitoClient) ([]types.AdEntry, error) {
	items := make([]types.AdEntry, 0)

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

func getDetailsForPage(url string, client *infra.IncognitoClient) (types.AdEntryDetails, error) {
	itemDetails := types.AdEntryDetails{}

	status, body, error := client.GetURLDataWithRetries(url, headers)

	if status == ErrBadRequest.Error() {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"url":    url,
		}).Warning("when fetch entry details")
		return itemDetails, error
	}

	if strings.Contains(status, "404") {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"url":    url,
		}).Warning("when fetch entry details")
		return itemDetails, nil
	}

	if error != nil {
		return itemDetails, error
	}

	doc, err := goquery.NewDocumentFromReader(body)

	if err != nil {
		return itemDetails, fmt.Errorf("unable to get document from response %w", err)
	}

	description := doc.Find(entityDetailsQuery).Text()

	itemDetails.Description = description

	var k []string

	var v []string

	doc.Find(".ClassifiedDetailBasicDetails-listTerm").Each(func(i int, s *goquery.Selection) {
		k = append(k, strings.TrimSpace(s.Text()))
	})

	doc.Find(".ClassifiedDetailBasicDetails-listDefinition").Each(func(i int, s *goquery.Selection) {
		v = append(v, strings.TrimSpace(s.Text()))
	})

	for i, key := range k {
		itemDetails.Attributes = append(itemDetails.Attributes, types.BasicInfoEntry{Key: key, Value: v[i]})
	}

	return itemDetails, nil
}

func getItemsForPage(url string, page int, client *infra.IncognitoClient) ([]types.AdEntry, error) {
	urlWithPage := url
	if page > 1 {
		u, err := buildURL(url, page)
		urlWithPage = u
		if err != nil {
			return nil, ErrBadRequest
		}
	}

	status, body, error := client.GetURLDataWithRetries(urlWithPage, headers)
	// Find review items
	items := []types.AdEntry{}

	if status == ErrBadRequest.Error() {
		// TODO(sycomancy): handle 400 BAD REQUEST!!!
		logrus.WithFields(logrus.Fields{
			"status": status,
			"url":    urlWithPage,
		}).Warning("fetching from njuskalo")
		return nil, ErrBadRequest
	}

	if strings.Contains(status, "404") {
		logrus.WithFields(logrus.Fields{
			"status": status,
			"url":    urlWithPage,
		}).Warning("fetching from njuskalo")
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
		var description string

		itemId, idExists := s.Attr("data-href")
		titleEl := s.Find(entityTitleQuery)
		hrefAttr, exists := titleEl.Attr("href")
		price = s.Find(entityPriceQuery).Text()
		title = titleEl.Text()
		description = s.Find(descriptionQuery).Text()

		if idExists {
			id = itemId
		}

		if exists {
			// TODO(sycomancy): add this to connfig
			link = "https://njuskalo.hr" + hrefAttr
		}

		item := types.AdEntry{Id: id, Title: title, Link: link, Price: price, Description: description}
		items = append(items, item)
	})

	return items, nil
}

func buildURL(base string, page int) (string, error) {
	u, err := url.Parse(base)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Add("page", fmt.Sprint(page))
	u.RawQuery = q.Encode()
	return u.String(), nil
}
