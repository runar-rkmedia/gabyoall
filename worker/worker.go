package worker

import (
	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/queries"
)

type WorkThing struct{}

func (w WorkThing) Run(endpoint cmd.Endpoint, config cmd.Config, query queries.GraphQLQuery) chan cmd.RequestStat {
	jobCh := make(chan Job, config.RequestCount)
	resultCh := make(chan cmd.RequestStat, config.RequestCount)
	// Create work
	for w := 0; w < config.Concurrency; w++ {
		// FIXME: this probably requires a lot of memory
		go worker(w, resultCh, jobCh)
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
	return resultCh
}

type Job struct {
	config   *cmd.Config
	endpoint *cmd.Endpoint
	query    *queries.GraphQLQuery
}

func worker(id int, ch chan cmd.RequestStat, jobCh chan Job) {
	for job := range jobCh {
		_, stat, _ := job.endpoint.RunQuery(*job.query, job.config.OkStatusCodes)
		if stat.Response != nil {

			if job.config.ResponseData != false && stat.Response != nil && stat.Response["error"] == nil && stat.Response["data"] != nil {
				delete(stat.Response, "data") //stat.Response[data]
			}
		}
		ch <- stat
	}
}

type Work func()
