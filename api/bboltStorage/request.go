package bboltStorage

import (
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
		Config: &p.Config,
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
	return e, err
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
