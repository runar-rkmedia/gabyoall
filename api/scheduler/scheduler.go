package scheduler

import (
	"fmt"
	"sync"
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
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

	config := s.config
	if ep.Config != nil {
		if ep.Config.Concurrency != nil && *ep.Config.Concurrency > 0 {
			config.Concurrency = *ep.Config.Concurrency
		}
		if ep.Config.RequestCount != nil && *ep.Config.RequestCount > 0 {
			config.RequestCount = *ep.Config.RequestCount
		}
		if ep.Config.OkStatusCodes != nil {
			config.OkStatusCodes = *ep.Config.OkStatusCodes
		}
		if ep.Config.ResponseData != nil {
			config.ResponseData = *ep.Config.ResponseData
		}
		if ep.Config.Auth != nil {
			if ep.Config.Auth.Endpoint != "" {
				config.Auth.Endpoint = ep.Config.Auth.Endpoint
			}
			if ep.Config.Auth.EndpointType != "" {
				config.Auth.EndpointType = ep.Config.Auth.EndpointType
			}
			if ep.Config.Auth.HeaderKey != "" {
				config.Auth.HeaderKey = ep.Config.Auth.HeaderKey
			}
			if ep.Config.Auth.Kind != "" {
				config.Auth.Kind = ep.Config.Auth.Kind
			}
			if ep.Config.Auth.RedirectUri != "" {
				config.Auth.RedirectUri = ep.Config.Auth.RedirectUri
			}
			if ep.Config.Auth.Token != "" {
				config.Auth.Token = string(ep.Config.Auth.Token)
			}
			if len(ep.Config.Auth.Dynamic.Requests) != 0 {
				for _, v := range ep.Config.Auth.Dynamic.Requests {
					config.Auth.Dynamic.Requests = make([]cmd.DynamicRequest, len(ep.Config.Auth.Dynamic.Requests))
					if vv := cmd.DynamicRequest(v); true {
						config.Auth.Dynamic.Requests = append(config.Auth.Dynamic.Requests, vv)
					}
				}
			}
			if ep.Config.Auth.ClientID != "" {
				config.Auth.ClientID = ep.Config.Auth.ClientID
			}
			if ep.Config.Auth.ClientSecret != "" && !types.IsRedacted(string(ep.Config.Auth.ClientSecret)) {
				config.Auth.ClientSecret = string(ep.Config.Auth.ClientSecret)
			}
			if ep.Config.Auth.ImpersionationCredentials.Password != "" && !types.IsRedacted(string(ep.Config.Auth.ImpersionationCredentials.Password)) {
				config.Auth.ImpersionationCredentials.Password = string(ep.Config.Auth.ImpersionationCredentials.Password)
			}
			if ep.Config.Auth.ImpersionationCredentials.UserIDToImpersonate != "" {
				config.Auth.ImpersionationCredentials.UserIDToImpersonate = ep.Config.Auth.ImpersionationCredentials.UserIDToImpersonate
			}
			if ep.Config.Auth.ImpersionationCredentials.UserNameToImpersonate != "" {
				config.Auth.ImpersionationCredentials.UserNameToImpersonate = ep.Config.Auth.ImpersionationCredentials.UserNameToImpersonate
			}
			if ep.Config.Auth.ImpersionationCredentials.Username != "" {
				config.Auth.ImpersionationCredentials.Username = ep.Config.Auth.ImpersionationCredentials.Username
			}
		}
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
	out, _ := cmd.NewOutput(l, "", config.Url, rq.Request, nil)
	startTime := time.Now()
	print := printer.NewPrinter(
		*config,
		nil,
		rq.OperationName,
		out,
		startTime,
	)
	print.Animate()
	config.Concurrency = 30
	config.RequestCount = 1000

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
	ch := wt.Run(endpoint, *config, rq.Request)
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
