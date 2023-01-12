package njuskalo

import (
	"fmt"
	"io"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/types"
)

var hostname = "https://njuskalo.hr/nekretnine/"

const entityListItemsQuery = ".EntityList--Standard > .EntityList-items > .EntityList-item.EntityList-item--Regular"
const entityTitleQuery = ".entity-title a"
const entityPriceQuery = ".price-item > .price"
const descriptionQuery = ".entity-description"

var leftDelimiter = "ina:"
var rightDelimiter = "m2"

var sizeRx = regexp.MustCompile(`(?s)` + regexp.QuoteMeta(leftDelimiter) + `(.*?)` + regexp.QuoteMeta(rightDelimiter))

var Headers = map[string]string{"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"}

var ErrBadRequest = fmt.Errorf("400 Bad Request")

func FetchAds(url string, client *infra.IncognitoClient) ([]types.AdEntry, error) {
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

func GetUrlForPage(url string, page int) (string, error) {
	urlWithPage := url
	if page != 0 {
		u, err := buildURL(url, page)
		urlWithPage = u
		if err != nil {
			return "", ErrBadRequest
		}
	}
	return urlWithPage, nil
}

func getPageReader(url string, page int, client *infra.IncognitoClient) (io.ReadCloser, error) {
	urlWithPage := url
	if page > 1 {
		u, err := buildURL(url, page)
		urlWithPage = u
		if err != nil {
			return nil, ErrBadRequest
		}
	}

	status, body, error := client.GetURLDataWithRetries(urlWithPage, Headers)
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
		return nil, error
	}

	if error != nil {
		return nil, error
	}
	return body, nil
}

func ParsePage(html string) ([]types.AdEntry, error) {
	reader := strings.NewReader(html)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		return nil, fmt.Errorf("unable to get document from response %w", err)
	}

	items := []types.AdEntry{}

	doc.Find(entityListItemsQuery).Each(func(i int, s *goquery.Selection) {
		var id string
		var title string
		var link string
		var price string
		var size string

		itemId, idExists := s.Attr("data-href")
		titleEl := s.Find(entityTitleQuery)
		hrefAttr, exists := titleEl.Attr("href")
		price = s.Find(entityPriceQuery).Text()
		title = titleEl.Text()

		description := s.Find(descriptionQuery).Text()
		matches := sizeRx.FindAllStringSubmatch(description, -1)
		size = strings.TrimSpace(matches[0][1])

		if idExists {
			id = itemId
		}

		if exists {
			// TODO(sycomancy): add this to connfig
			link = "https://njuskalo.hr" + hrefAttr
		}

		item := types.AdEntry{Id: id, Title: title, Link: link, Price: price, Size: size}
		items = append(items, item)
	})

	return items, nil
}

func getItemsForPage(url string, page int, client *infra.IncognitoClient) ([]types.AdEntry, error) {
	items := []types.AdEntry{}
	body, err := getPageReader(url, page, client)
	if err != nil {
		return nil, err
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
			// TODO(sycomancy): add this to connfig
			link = "https://njuskalo.hr" + hrefAttr
		}

		item := types.AdEntry{Id: id, Title: title, Link: link, Price: price}
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
