package njuskalo

import (
	"fmt"
	"testing"

	"github.com/sycomancy/glasnik/infra"
)

var urls = []string{"", "", ""}

func TestFetchItemsForNjuska(t *testing.T) {
	client := infra.NewIncognitoClient(nil)

	for _, url := range urls {
		result, err := FetchAds(url, client)

		if err != nil {
			t.Error(err)
		} else {
			count := len(result)
			fmt.Printf("Got %d results for url: %s \n", count, url)
		}
	}

}
