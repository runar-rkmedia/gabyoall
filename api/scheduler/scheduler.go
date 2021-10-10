package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/requests"
	"github.com/runar-rkmedia/gabyoall/worker"
)

type Scheduler struct {
	l         logger.AppLogger
	db        types.Storage
	interval  time.Duration
	config    *cmd.Config
	isRunning bool
	sync.Mutex
}

func nilDurationStr(t *time.Time) string {
	if t == nil {
		return ""
	}
	now := time.Now()
	return now.Sub(*t).String()
}

func (s *Scheduler) Run() {
	ticker := time.NewTicker(s.interval)
	quit := make(chan struct{})
	debug := s.l.HasDebug()
	s.l.Info().Dur("interval", s.interval).Msg("Starting scheduler at interval")
	go func() {
		for {
			select {
			case <-ticker.C:
				s.Lock()
				if s.isRunning {
					if debug {
						s.l.Debug().Msg("Already running")
					}
					continue
				}
				s.Unlock()

				schedule, err := s.db.Schedules()
				if err != nil {
					continue
				}
				for _, v := range schedule {
					if !v.ShouldRun() {
						continue
					}
					err := s.RunSchedule(v)
					v.Schedule.LastError = err
					now := time.Now()
					v.LastRun = &now
					_, err = s.db.UpdateSchedule(v.ID, v.Schedule)
					if err != nil {
						s.l.Error().Err(err).Msg("There was a problem updating the schedule")
					}

				}
			case <-quit:
				ticker.Stop()
				return

			}
		}
	}()
}

func (s *Scheduler) RunSchedule(v types.ScheduleEntity) error {
	debug := s.l.HasDebug()
	l := logger.With(s.l.With().
		Str("RequestID", v.RequestID).
		Str("EndpointID", v.EndpointID).
		Logger())

	// TODO: Check further to see if we really are scheduled to run.
	wt := worker.WorkThing{}
	ep, err := s.db.Endpoint(v.EndpointID)
	if err != nil {
		l.Error().Str("ID", v.EndpointID).Err(err).Msg("Failed to get endpoint for id")
		return err
	}
	request, err := s.db.Request(v.RequestID)
	if err != nil {
		l.Error().Str("ID", v.RequestID).Err(err).Msg("Failed to get request for id")
		return err
	}

	config := s.config
	if config.Concurrency == 0 {
		s.l.Error().Msg("Concurrency must be positive")
		return fmt.Errorf("Concurrency must be positive")
	}
	if config.RequestCount == 0 {
		s.l.Error().Msg("RequestCount must be positive")
		return fmt.Errorf("RequestCount must be positive")
	}
	if debug {
		l.Debug().
			Int("concurrency", config.Concurrency).
			Int("requestCount", config.RequestCount).
			Msg("Starting request")
	}
	endpoint := requests.NewEndpoint(s.l, ep.Url)
	// TODO: get config from db-entities and merge
	ch := wt.Run(endpoint, *config, request.Request)
	for i := 0; i < config.RequestCount; i++ {
		stat := <-ch
		fmt.Println(stat)
	}
	return nil
}

func NewScheduler(l logger.AppLogger, db types.Storage, config *cmd.Config) Scheduler {
	s := Scheduler{
		l:        l,
		interval: 500 * time.Millisecond,
		db:       db,
		config:   config,
	}
	return s
}
