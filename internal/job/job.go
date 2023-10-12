package job

import (
	"fmt"

	"github.com/sycomancy/glasnik/internal/infra"
	"github.com/sycomancy/glasnik/internal/njuskalo"
	"github.com/sycomancy/glasnik/internal/types"
)

type Job struct {
	Url    string
	client *infra.IncognitoClient
}

func NewJob(url string) (*Job, error) {
	job := &Job{Url: url}
	job.client = infra.NewIncognitoClient(nil)
	return job, nil
}

func (j *Job) Start() error {
	itemCh := make(chan types.AdEntry)
	go njuskalo.FetchEntry(j.Url, itemCh, j.client)

	for item := range itemCh {
		fmt.Println(item)
	}
	// entries, err := njuskalo.Fetch(j.Url, &j.client)
	// if err != nil {
	// 	return err
	// }

	// for _, entry := range entries {
	// 	fmt.Println(entry.Id)
	// }

	return nil
}

// func NewFetchJob(urls []string, interval int) (*FetchJob, error) {
// 	job := &FetchJob{Urls: urls}
// 	return job, nil
// }

// func (f *FetchJob) start() {
// 	ticker := time.NewTicker(5 * time.Second)
// 	quit := make(chan struct{})
// 	go func() {
// 		for {
// 			select {
// 			case <-ticker.C:
// 				// do stuff
// 			case <-quit:
// 				ticker.Stop()
// 				return
// 			}
// 		}
// 	}()
// }
