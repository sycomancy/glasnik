package analytics

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/sycomancy/glasnik/service"
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "fetchSchedule",
})

// type LocationDetails struct {
// 	Id    string `json:"id,omitempty"`
// 	Title string `json:"title,omitempty"`
// }

type AnalyticResult struct {
	JobID string `json:"jobId,omitempty"`
	Date  string `json:"date,omitempty"`
	// Location         LocationDetails `json:"location,omitempty"`
	Count            int `json:"count,omitempty"`
	MaxPrice         int `json:"maxPrice,omitempty"`
	MinPrice         int `json:"minPrice,omitempty"`
	AvgPrice         int `json:"avgPrice,omitempty"`
	StdPrice         int `json:"stdPrice,omitempty"`
	AvgPricePerMeter int `json:"avgPricePerMeter,omitempty"`
}

type Generator struct {
	storer *service.Storer
}

func NewGenerator(storer *service.Storer) *Generator {
	return &Generator{storer: storer}
}

func (g *Generator) Process(jobID string) {
	// get all results for job
	results, err := g.storer.ResultsForJob(jobID)
	if err != nil {
		logg.Fatalf("failed to fetch results for %s %w", jobID, err)
	}

	forLocResult := make(chan string)
	for _, r := range results {
		go g.calculatePerLocation(r, forLocResult)
	}

	for res := range forLocResult {
		fmt.Printf("Received result for %s \n", res)
	}
}

func (g *Generator) calculatePerLocation(loc *service.LocationPageResult, result chan<- string) {
	time.Sleep(time.Second * 2)
	fmt.Println("Hereeeee", loc)
	result <- loc.Location.Id
}
