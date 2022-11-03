package njuskalo

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"

	"github.com/damjackk/njufetch/infra"
)

const entityListItemsQuery = ".EntityList--Standard > .EntityList-items > .EntityList-item.EntityList-item--Regular"
const entityTitleQuery = ".entity-title a"
const entityPriceQuery = ".price-item > .price"

var headers = map[string]string{"user-agent": "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/104.0.0.0 Safari/537.36"}

var ErrBadRequest = fmt.Errorf("400 Bad Request")

type EntityItem struct {
	Id    string
	Title string
	Link  string
	Price string
}

func Fetch(url string, client *infra.IncognitoClient) ([]EntityItem, error) {
	items := make([]EntityItem, 0)

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
			fmt.Printf("Fetched: %d%s%d\n", len(pageItems), " for page: ", page)
			page += 1
		}
	}

	if len(items) != 0 {
		page += 1
	}

	return items, nil
}

func getItemsForPage(url string, page int, client *infra.IncognitoClient) ([]EntityItem, error) {
	urlWithPage := url

	if page != 1 {
		urlWithPage = fmt.Sprintf("%s%s%d", url, "&page=", page)
	}

	status, body, error := client.GetURLDataWithRetries(urlWithPage, headers)
	fmt.Println("getItemsForPage ", status)

	// Find review items
	items := []EntityItem{}

	if status == ErrBadRequest.Error() {
		// TODO(sycomancy): handle 400 BAD REQUEST!!!
		fmt.Printf("WE got 400 for %s \n", status)
		return nil, ErrBadRequest
	}

	if strings.Contains(status, "404") {
		fmt.Printf("WE got 400 for url %s", status)
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

		if idExists {
			id = itemId
		}

		if exists {
			link = hrefAttr
		}

		item := EntityItem{Id: id, Title: title, Link: link, Price: price}
		items = append(items, item)
	})

	return items, nil
}
