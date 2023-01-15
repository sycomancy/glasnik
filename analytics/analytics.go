package analytics

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"github.com/sycomancy/glasnik/service"
)

var logg = logrus.WithFields(logrus.Fields{
	"ctx": "fetchSchedule",
})

type LocationDetails struct {
	Id    string `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type AnalyticResult struct {
	JobID            string          `json:"jobId,omitempty"`
	Date             string          `json:"date,omitempty"`
	Location         LocationDetails `json:"location,omitempty"`
	Count            int             `json:"count,omitempty"`
	MaxPrice         int             `json:"maxPrice,omitempty"`
	MinPrice         int             `json:"minPrice,omitempty"`
	AvgPrice         int             `json:"avgPrice,omitempty"`
	StdPrice         int             `json:"stdPrice,omitempty"`
	AvgPricePerMeter int             `json:"avgPricePerMeter,omitempty"`
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

	for _, r := range results {
		fmt.Println(r.Completed)
	}
}

func (g *Generator) CalculatePerLocation()
