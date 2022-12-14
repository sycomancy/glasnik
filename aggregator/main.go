package aggregator

import (
	"github.com/sirupsen/logrus"

	"github.com/sycomancy/glasnik/infra"
	"github.com/sycomancy/glasnik/njuskalo"
	"github.com/sycomancy/glasnik/types"
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "aggregator",
})

var filter = "https://www.njuskalo.hr/prodaja-stanova/zagreb"

type aggregator struct {
	client *infra.IncognitoClient
}

func (a *aggregator) FetchAndPersist(filters []string) {
	results, err := njuskalo.Fetch(filter, a.client)
	if err != nil {
		logg.Error(err)
	}

	a.persist(results)
	logg.Info("got: ", len(results), " results for filter: ", filter)
}

func (a *aggregator) persist(results []types.AdEntry) {
	entries := make([]interface{}, len(results))
	for i := range results {
		entries[i] = results[i]
	}
	_ = infra.InsertDocuments("results", entries)
}

func Run() {
	aggregator := &aggregator{client: infra.NewIncognitoClient(nil)}
	aggregator.FetchAndPersist([]string{filter})
}
