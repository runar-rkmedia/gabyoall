package bboltStorage

import (
	"time"

	"github.com/imdario/mergo"
	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/requests"
	bolt "go.etcd.io/bbolt"
)

func (s *BBolter) Request(id string) (e types.RequestEntity, err error) {
	err = s.GetItem(BucketRequests, id, &e)
	return
}

func (s *BBolter) CreateRequest(p types.RequestPayload) (types.RequestEntity, error) {
	e := types.RequestEntity{
		Request: requests.Request{
			Body:          p.Body,
			Query:         p.Query,
			Variables:     p.Variables,
			Headers:       p.Headers,
			OperationName: p.OperationName,
			Method:        p.Method,
		},
		Config: p.Config,
		Entity: s.NewEntity(),
	}

	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketRequests)
		bytes, err := s.Marshal(e)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(e.ID), bytes)
	})
	if err == nil {
		s.PublishChange(PubTypeRequest, PubVerbCreate, e)
	}
	return e, err
}

func (s *BBolter) UpdateRequest(id string, p types.RequestPayload) (types.RequestEntity, error) {
	var j types.RequestEntity
	if id == "" {
		return j, ErrMissingIdArg
	}
	request := types.RequestEntity{
		Request: requests.Request{
			Body:          p.Body,
			Query:         p.Query,
			Variables:     p.Variables,
			Headers:       p.Headers,
			OperationName: p.OperationName,
			Method:        p.Method,
		},
		Config: p.Config,
	}
	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketRequests)
		b := bucket.Get([]byte(id))
		if len(b) == 0 {
			return ErrNotFound
		}

		err := s.Unmarshal(b, &j)
		if err != nil {
			return err
		}
		err = mergo.Merge(&j.Request, request.Request, mergo.WithOverride)
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
		s.l.Error().Err(err).Msg("Failed during UpdateRequest")
	}
	return j, err
}

func (s *BBolter) Requests() (es map[string]types.RequestEntity, err error) {
	es = map[string]types.RequestEntity{}
	err = s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(BucketRequests)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e types.RequestEntity
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

func (s *BBolter) softDeleteRequest(id string, delete *bool) (j types.RequestEntity, err error) {
	err = s.updater(id, BucketRequests, func(b []byte) ([]byte, error) {
		if err := s.Unmarshal(b, &j); err != nil {
			return nil, err
		}
		now := time.Now()
		if delete == nil {
			if j.Deleted == nil {
				j.Deleted = &now
			} else {
				j.Deleted = nil
			}
		} else if *delete == false {
			j.Deleted = nil
		} else if *delete == true {
			j.Deleted = &now
		}
		j.UpdatedAt = &now
		return s.Marshal(j)
	})
	if err == nil {
		s.PublishChange(PubTypeRequest, PubVerbSoftDelete, j)
	}

	return
}

func (s *BBolter) SoftDeleteRequest(id string) (j types.RequestEntity, err error) {
	return s.softDeleteRequest(id, nil)
}
