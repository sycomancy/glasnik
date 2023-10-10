package main

import "time"

type WatchJob struct {
	Url      string
	Interval int
}

func (j *WatchJob) Start() {
	ticker := time.NewTicker(5 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				// do stuff
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}

func NewWatchJob(url string, interval int) (*WatchJob, error) {
	job := &WatchJob{Url: url, Interval: interval}
	return job, nil
}

type Watcher struct {
	jobs []WatchJob
}
