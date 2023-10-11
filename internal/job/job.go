package job

import (
	"fmt"

	"github.com/sycomancy/glasnik/internal/infra"
)

type Job struct {
	Urls   []string
	client infra.IncognitoClient
}

func NewJob(urls []string) (*Job, error) {
	job := &Job{Urls: urls}
	job.client = *infra.NewIncognitoClient(nil)
	return job, nil
}

func (j *Job) Start() {
	// result, err := njuskalo.Fetch()
	for _, url := range j.Urls {
		fmt.Println(url)
	}
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
