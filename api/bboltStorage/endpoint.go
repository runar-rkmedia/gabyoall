package bboltStorage

import (
	"net/http"
	"time"

	"github.com/imdario/mergo"
	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/requests"
	"github.com/runar-rkmedia/gabyoall/utils"
	bolt "go.etcd.io/bbolt"
)

func (s *BBolter) Endpoint(id string) (e types.EndpointEntity, err error) {
	err = s.GetItem(BucketEndpoints, id, &e)
	return
}
func (s *BBolter) CreateEndpoint(p types.EndpointPayload) (types.EndpointEntity, error) {
	entity, _ := ForceNewEntity()
	e := types.EndpointEntity{
		Endpoint: types.Endpoint{
			Endpoint: requests.Endpoint{
				Url:     p.Url,
				Headers: http.Header(p.Headers),
			},
			Config: p.Config,
		},
		Entity: entity,
	}

	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketEndpoints)
		bytes, err := s.Marshal(e)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(e.ID), bytes)
	})
	if err == nil {
		s.PublishChange(PubTypeEndpoint, PubVerbCreate, entity)
	}
	return e, err
}

func (s *BBolter) softDeleteEndpoint(id string, delete *bool) (j types.EndpointEntity, err error) {
	err = s.updater(id, BucketEndpoints, func(b []byte) ([]byte, error) {
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
		s.PublishChange(PubTypeEndpoint, PubVerbSoftDelete, j)
	}

	return
}
func (s *BBolter) SoftDeleteEndpoint(id string) (j types.EndpointEntity, err error) {
	return s.softDeleteEndpoint(id, nil)
}

func (s *BBolter) UpdateEndpoint(id string, p types.EndpointPayload) (types.EndpointEntity, error) {
	var j types.EndpointEntity
	if id == "" {
		return j, ErrMissingIdArg
	}
	endpoint := types.Endpoint{
		Endpoint: requests.Endpoint{
			Url:     p.Url,
			Headers: p.Headers,
		},
		Config: p.Config,
	}
	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketEndpoints)
		b := bucket.Get([]byte(id))
		if len(b) == 0 {
			return ErrNotFound
		}

		err := s.Unmarshal(b, &j)
		if err != nil {
			return err
		}
		err = mergo.Merge(&j.Endpoint, endpoint, mergo.WithOverride)
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
		s.l.Error().Err(err).Msg("Failed during UpdateEndpoint")
	}
	return j, err
}

func (s *BBolter) Endpoints() (es map[string]types.EndpointEntity, err error) {
	es = map[string]types.EndpointEntity{}
	err = s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(BucketEndpoints)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e types.EndpointEntity
			err := s.Unmarshal(v, &e)
			if err != nil {
				return err
			}
			es[string(k)] = e
		}

		return nil
	})
	if err != nil {
		s.l.Error().Err(err).Msg("failed to lookup endpoints")
	}

	return es, err
}

// Returns an entity for use by database, with id set and createdAt to current time.
// It is guaranteeed to be created correctly, if if it errors.
// The error should be logged.
func ForceNewEntity() (types.Entity, error) {
	id, err := utils.ForceCreateUniqueId()

	return types.Entity{
		ID:        id,
		CreatedAt: time.Now(),
	}, err
}
