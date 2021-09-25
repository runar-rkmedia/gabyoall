package worker

import (
	"math/rand"
	"time"

	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/queries"
)

type WorkThing struct{}

func (w WorkThing) Run(endpoint cmd.GraphQlEndpoint, config cmd.Config, query queries.GraphQLQuery) (chan cmd.RequestStat, chan struct{}) {
	hasWorkChan := make(chan struct{}, config.Concurrency)
	workChan := make(chan Work, config.Concurrency)
	ch := make(chan cmd.RequestStat)
	go func() {
		for {
			select {
			case work := <-workChan:
				go work()
				hasWorkChan <- struct{}{}

			}

		}
	}()
	go func() {
		for j := 0; j < config.RequestCount; j++ {

			workChan <- func(i int) Work {

				return func() {
					if config.Mock {
						// TODO: replace with a mocked http-client-interface
						stat := cmd.NewStat()
						time.Sleep(time.Millisecond * time.Duration(rand.Int63n(80)+1))
						errorType := cmd.Unknwon
						n := rand.Intn(7)
						switch n {
						case 1:
							errorType = cmd.NonOK
						case 2:
							errorType = cmd.ServerTestError
						case 3:
							errorType = "RandomErr"
						case 4:
							errorType = "OtherErr"
						case 6:
							errorType = "MadeUpError"

						}
						ch <- stat.End(errorType, nil)
						return
					}
					_, stat, _ := endpoint.RunQuery(query, config.OkStatusCodes)
					if stat.Response != nil {

						if config.ResponseData != false && stat.Response != nil && stat.Response["error"] == nil && stat.Response["data"] != nil {
							delete(stat.Response, "data") //stat.Response[data]
						}
					}
					ch <- stat

				}
			}(j)
		}
	}()
	return ch, hasWorkChan
}

type Work func()
