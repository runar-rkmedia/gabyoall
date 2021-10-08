package bboltStorage

import (
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/logger"
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
func (s *BBolter) NewEntity() types.Entity {
	// ForceNewEntity may return an error, but it guarantees the the Entity is still usable.
	// The error should be logged, though.
	e, err := ForceNewEntity()
	if err != nil {
		s.l.Error().Err(err).Str("id", e.ID).Msg("An error occured while creating entity. ")
	}
	return e
}

type BBolter struct {
	*bolt.DB
	l logger.AppLogger
	Marshaller
}

var (
	BucketEndpoints = []byte("endpoints")
	BucketRequests  = []byte("requests")
)
