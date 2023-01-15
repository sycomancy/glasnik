package njuskalo

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/sycomancy/glasnik/types"
)

const entityListItemsQuery = ".EntityList--Standard > .EntityList-items > .EntityList-item.EntityList-item--Regular"
const entityTitleQuery = ".entity-title a"
const entityPriceQuery = ".price-item > .price"
const descriptionQuery = ".entity-description"

var Headers = map[string]string{"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"}

var ErrBadRequest = fmt.Errorf("400 Bad Request")

// ------ Parsersssssssssssssssss --------------------
var leftDelimiter = "ina:"
var rightDelimiter = "m2"
var sizeRx = regexp.MustCompile(`(?s)` + regexp.QuoteMeta(leftDelimiter) + `(.*?)` + regexp.QuoteMeta(rightDelimiter))

func parseSize(txt string) string {
	matches := sizeRx.FindAllStringSubmatch(txt, -1)
	return strings.TrimSpace(matches[0][1])
}

var leftPriceDelimiter = ""
var rightPriceDelimiter = "€"
var priceRx = regexp.MustCompile(`(?s)` + regexp.QuoteMeta(leftPriceDelimiter) + `(.*?)` + regexp.QuoteMeta(rightPriceDelimiter))

func parsePrice(txt string) string {
	matches := priceRx.FindAllStringSubmatch(txt, -1)
	return strings.TrimSpace(matches[0][1])
}

// ---------------------------------------------------
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
		price = parsePrice(s.Find(entityPriceQuery).Text())
		title = titleEl.Text()
		description := s.Find(descriptionQuery).Text()
		size = parseSize(description)

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
