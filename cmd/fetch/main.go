package main

import (
	"fmt"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
)

var urls = []string{"https://www.njuskalo.hr/prodaja-stanova?geo%5BlocationIds%5D=2698"}

func main() {
	infra.LoadConfig()

	client := infra.NewIncognitoClient(nil)
	fmt.Print(client)
	for _, url := range urls {
		result, err := njuskalo.Fetch(url, client)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Printf("Got %d results for url: %s \n", len(result), url)
			for _, r := range result {
				fmt.Println(r.Title, r.Id, r.Price)
			}
		}
	}
}
