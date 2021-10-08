package bboltStorage

import (
	"net/http"
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/requests"
	bolt "go.etcd.io/bbolt"
)

func (s *BBolter) Endpoint(id string) (e types.EndpointEntity, err error) {
	err = s.GetItem(BucketEndpoints, id, &e)
	return
}
func (s *BBolter) CreateEndpoint(p types.EndpointPayload) (types.EndpointEntity, error) {
	id, _ := ForceCreateUniqueId()
	now := time.Now()
	e := types.EndpointEntity{
		Endpoint: requests.Endpoint{
			Url:     p.Url,
			Headers: http.Header(p.Headers),
		},
		Entity: types.Entity{
			ID:        id,
			CreatedAt: now,
		},
	}

	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketEndpoints)
		bytes, err := s.Marshal(e)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(e.ID), bytes)
	})
	return e, err
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
