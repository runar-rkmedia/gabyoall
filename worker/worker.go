package worker

import (
	"time"

	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/requests"
)

type WorkThing struct{}

func (w WorkThing) Run(endpoint requests.Endpoint, config cmd.Config, query requests.Request) (chan requests.RequestStat, chan Job) {
	jobCh := make(chan Job, config.RequestCount)
	resultCh := make(chan requests.RequestStat, config.RequestCount)
	startTime := time.Now()
	// Create work
	for w := 0; w < config.Concurrency; w++ {
		// FIXME: this probably requires a lot of memory
		go worker(startTime, w, resultCh, jobCh)
	}

	// Create
	for j := 0; j < config.RequestCount; j++ {
		job := Job{
			config:   &config,
			endpoint: &endpoint,
			query:    &query,
		}
		go func(job Job) {
			jobCh <- job
		}(job)
	}
	return resultCh, jobCh
}

type Job struct {
	config   *cmd.Config
	endpoint *requests.Endpoint
	query    *requests.Request
}

func worker(startTime time.Time, id int, ch chan requests.RequestStat, jobCh chan Job) {
	for job := range jobCh {
		_, stat, _ := job.endpoint.RunQuery(startTime, *job.query, job.config.OkStatusCodes)
		ch <- stat
	}
}

type Work func()
