package bboltStorage

import (
	"fmt"
	"time"

	"github.com/imdario/mergo"
	"github.com/runar-rkmedia/gabyoall/api/types"
	bolt "go.etcd.io/bbolt"
)

func (s *BBolter) Schedule(id string) (e types.ScheduleEntity, err error) {
	err = s.GetItem(BucketSchedules, id, &e)
	return
}

// Used to nullcheck a time. Useful when dealing with non-null field
func isTimeSet(d time.Time) bool {
	minDate := time.Time{}
	minDate.Add(1)
	return d.After(minDate)
}

func (s *BBolter) UpdateSchedule(id string, p types.Schedule) (types.ScheduleEntity, error) {
	var j types.ScheduleEntity
	if id == "" {
		return j, ErrMissingIdArg
	}
	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSchedules)
		b := bucket.Get([]byte(id))
		if len(b) == 0 {
			return ErrNotFound
		}

		err := s.Unmarshal(b, &j)
		if err != nil {
			return err
		}
		err = mergo.Merge(&j.Schedule, p, mergo.WithOverride)
		if err != nil {
			return err
		}
		if p.Config != nil {
			if j.Config == nil {
				j.Config = &types.Config{}
			}
			err = mergo.Merge(j.Config, p.Config, mergo.WithOverride)
			if err != nil {
				return err
			}
		}
		now := time.Now()
		j.UpdatedAt = &now
		bytes, err := s.Marshal(j)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(id), bytes)

	})
	if err != nil {
		s.l.Error().Err(err).Msg("Failed during UpdateSchedule")
	}
	return j, err
}
func (s *BBolter) CreateSchedule(p types.SchedulePayload) (types.ScheduleEntity, error) {
	e := types.ScheduleEntity{
		Entity: s.NewEntity(),
		Config: p.Config,
		Schedule: types.Schedule{
			Dates:           []time.Time{},
			SchedulePayload: p,
		},
	}
	ep, err := s.Endpoint(e.EndpointID)
	if err != nil {
		return e, err
	}
	if ep.ID == "" {
		return e, fmt.Errorf("EndpointID is missing")
	}
	rq, err := s.Request(e.RequestID)
	if err != nil {
		return e, err
	}
	if rq.ID == "" {
		return e, fmt.Errorf("RequestID is missing")
	}
	// TODO: create Dates for the next x / year

	err = s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketSchedules)
		bytes, err := s.Marshal(e)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(e.ID), bytes)
	})
	return e, err
}

func (s *BBolter) Schedules() (es map[string]types.ScheduleEntity, err error) {
	es = map[string]types.ScheduleEntity{}
	err = s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(BucketSchedules)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e types.ScheduleEntity
			err := s.Unmarshal(v, &e)
			if err != nil {
				return err
			}
			es[string(k)] = e
		}

		return nil
	})
	if err != nil {
		s.l.Error().Err(err).Msg("failed to lookup requests")
	}

	return es, err
}
