package scheduler

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/auth"
	"github.com/runar-rkmedia/gabyoall/cmd"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/printer"
	"github.com/runar-rkmedia/gabyoall/requests"
	"github.com/runar-rkmedia/gabyoall/utils"
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
					if err != nil {
						v.Schedule.LastError = err.Error()
					} else {
						// TODO: fix this hack
						v.Schedule.LastError = "__CLEAR__"
					}
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
	l := logger.With(s.l.With().
		Str("RequestID", v.RequestID).
		Str("EndpointID", v.EndpointID).
		Str("EndpointLabel", v.Label).
		Logger())

	// TODO: Check further to see if we really are scheduled to run.
	ep, err := s.db.Endpoint(v.EndpointID)
	if err != nil {
		l.Error().Str("ID", v.EndpointID).Err(err).Msg("Failed to get endpoint for id")
		return err
	}
	rq, err := s.db.Request(v.RequestID)
	if err != nil {
		l.Error().Str("ID", v.RequestID).Err(err).Msg("Failed to get request for id")
		return err
	}

	config := cmd.Config{}
	// order matters
	configs := []*types.Config{
		ep.Config,
		rq.Config,
		v.Config,
	}
	for i := 0; i < len(configs); i++ {
		if configs[i] == nil {
			continue
		}
		config = configs[i].MergeInto(config)
	}
	if config.Concurrency == 0 {
		s.l.Error().Msg("Concurrency must be positive")
		return fmt.Errorf("Concurrency must be positive")
	}
	if config.RequestCount == 0 {
		s.l.Error().Msg("RequestCount must be positive")
		return fmt.Errorf("RequestCount must be positive")
	}
	endpoint := requests.NewEndpoint(s.l, ep.Url)
	var token string
	// TODO: renew the tokenPayload as needed
	// var tokenPayload *auth.TokenPayload
	var validityStringer printer.ValidityStringer
	if token == "" {
		err, token, _, validityStringer = auth.Retrieve(l, config.Auth)
		if err != nil {
			return fmt.Errorf("failed to perform authentication: %w", err)
		}
	}
	if token != "" {
		authPrefix := ""
		if strings.ToLower(config.Auth.Kind) == "bearer" {
			authPrefix = "Bearer "
		}
		if config.Auth.HeaderKey == "" {
			config.Auth.HeaderKey = "Authorization"
		}
		endpoint.Headers.Add(config.Auth.HeaderKey, authPrefix+token)
	}
	out, _ := cmd.NewOutput(l, "", config.Url, rq.Request, nil)
	startTime := time.Now()
	print := printer.NewPrinter(
		config,
		validityStringer,
		rq.OperationName,
		out,
		startTime,
	)
	print.Animate()

	wt := worker.WorkThing{}
	l = logger.With(s.l.With().
		Int("Concurrency", config.Concurrency).
		Int("Request-Count", config.RequestCount).
		Str("OperationName", rq.OperationName).
		Str("URL", ep.Url).
		Str("Method", rq.Method).
		Logger())
	l.Info().
		Msg("Starting scheduled request")
	startedAt := time.Now()
	runId, warn := utils.ForceCreateUniqueId()
	if warn != nil {
		l.Warn().Err(warn).Str("ID", runId).Msg("Failed to create id, used fallback-method instead")
	}
	ch, jobCh := wt.Run(endpoint, config, rq.Request)
	defer close(jobCh)
	defer close(ch)
	successes := 0
	stats := requests.NewCompactRequestStatistics()
	lastSave := time.Now()
	debug := s.l.HasDebug()
	didSave := false
	for i := 0; i < config.RequestCount; i++ {
		stat := <-ch
		if stat.ErrorType == "" {
			successes++
		}
		stats.AddStat(stat)
		now := time.Now()
		if i == config.RequestCount || now.Sub(lastSave) > time.Second {
			lastSave = now
			stats.Calculate()
			if debug {
				s.l.Debug().Msg("Saving")
			}
			if didSave {
				s.db.UpdateCompactStats(runId, startedAt, stats)
			} else {
				s.db.CreateCompactStats(runId, startedAt, stats)
			}
			didSave = true
		}
	}
	l.Info().
		Msg("Completed scheduled request")
	if didSave {
		s.db.UpdateCompactStats(runId, startedAt, stats)
	} else {
		s.db.CreateCompactStats(runId, startedAt, stats)
	}
	return nil
}

func NewScheduler(l logger.AppLogger, db types.Storage, config *cmd.Config) *Scheduler {
	s := Scheduler{
		l:        l,
		interval: 500 * time.Millisecond,
		db:       db,
		config:   config,
	}
	return &s
}
