package bboltStorage

import (
	"errors"
	"os"
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
	"github.com/runar-rkmedia/gabyoall/logger"
	bolt "go.etcd.io/bbolt"
)

var (
	ErrMissingIdArg = errors.New("Missing id as argument")
	ErrNotFound     = errors.New("Not found")
)

type PubSubPublisher interface {
	Publish(kind, variant string, contents interface{})
}

// Caller must call close when ending
func NewBbolt(l logger.AppLogger, path string, pubsub PubSubPublisher) (bb BBolter, err error) {

	bb.l = l
	db, err := bolt.Open(path, 0666, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		return
	}
	bb.DB = db
	bb.pubsub = pubsub
	bb.Marshaller = Gob{}
	err = bb.Update(func(t *bolt.Tx) error {
		buckets := [][]byte{BucketEndpoints, BucketRequests, BucketSchedules, BucketStats}
		for i := 0; i < len(buckets); i++ {
			_, err := t.CreateBucketIfNotExists(buckets[i])
			if err != nil {
				return err

			}
		}
		return nil
	})
	return
}

func (s *BBolter) PublishChange(kind PubType, variant PubVerb, contents interface{}) {
	if s.pubsub == nil {
		return
	}
	s.pubsub.Publish(string(kind), string(variant), contents)
}

func (s *BBolter) GetItem(bucket []byte, id string, j interface{}) error {
	if id == "" {
		return ErrMissingIdArg
	}
	err := s.DB.View(func(t *bolt.Tx) error {
		bucket := t.Bucket(bucket)
		b := bucket.Get([]byte(id))
		if b == nil || len(b) == 0 {
			return ErrNotFound
		}
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

func (s *BBolter) Size() (int64, error) {
	s.l.Info().Interface("stats", s.Stats()).Msg("DB-stats")

	stat, err := os.Stat(s.Path())
	if err != nil {
		return 0, err
	}
	return int64(stat.Size()), err
}
func (s *BBolter) updater(id string, bucket []byte, f func(b []byte) ([]byte, error)) error {
	if id == "" {
		return ErrMissingIdArg
	}
	if bucket == nil {
		return ErrMissingIdArg
	}
	err := s.Update((func(tx *bolt.Tx) error {
		bucket := tx.Bucket(bucket)
		b := bucket.Get([]byte(id))
		if len(b) == 0 {
			return ErrNotFound
		}
		newBytes, err := f(b)
		if err != nil {
			return err
		}

		return bucket.Put([]byte(id), newBytes)
	}))

	return err
}

type BBolter struct {
	*bolt.DB
	pubsub PubSubPublisher
	l      logger.AppLogger
	Marshaller
}

var (
	BucketEndpoints = []byte("endpoints")
	BucketRequests  = []byte("requests")
	BucketSchedules = []byte("schedules")
	BucketStats     = []byte("stats")
)
