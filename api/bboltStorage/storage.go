package bboltStorage

import (
	"time"

	"net/http"

	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/logger"
	"github.com/runar-rkmedia/gabyoall/requests"
	bolt "go.etcd.io/bbolt"
)

// Caller must call close when ending
func NewBbolt(l logger.AppLogger, path string) (bb BBolter, err error) {

	bb.l = l
	db, err := bolt.Open(path, 0666, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return
	}
	bb.DB = db
	bb.Marshaller = Gob{}
	err = bb.Update(func(t *bolt.Tx) error {
		_, err := t.CreateBucketIfNotExists(BucketEndpoints)
		return err
	})
	return
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
func (s *BBolter) GetItem(bucket []byte, id string, j interface{}) error {
	err := s.DB.View(func(t *bolt.Tx) error {
		bucket := t.Bucket(bucket)
		b := bucket.Get([]byte(id))
		return s.Unmarshal(b, j)
	})
	if err != nil {
		s.l.Error().Err(err).Bytes("bucket", bucket).Str("id", id).Msg("Failed to lookup endpoint")
	}

	return err

}
func (s *BBolter) Endpoint(id string) (e types.EndpointEntity, err error) {
	err = s.GetItem(BucketEndpoints, id, &e)
	return
}
func (s *BBolter) CreateEndpoint(p types.EndpointPayload) (types.EndpointEntity, error) {
	id, _ := CreateUniqueId()
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

type BBolter struct {
	*bolt.DB
	l logger.AppLogger
	Marshaller
}

var (
	BucketEndpoints = []byte("1")
)
