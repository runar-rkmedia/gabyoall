package bboltStorage

import (
	"fmt"
	"time"

	"github.com/runar-rkmedia/gabyoall/api/types"
	bolt "go.etcd.io/bbolt"
)

func (s *BBolter) UpdateCompactStats(id string, createdAt time.Time, p types.StatPayload) error {
	now := time.Now()
	e := types.StatEntity{
		Entity: types.Entity{
			ID:        id,
			CreatedAt: createdAt,
			UpdatedAt: &now,
		},
		CompactRequestStatistics: p,
	}
	return s.writeCompactStats(id, e)
}

func (s *BBolter) CreateCompactStats(id string, createAt time.Time, p types.StatPayload) error {
	now := time.Now()
	e := types.StatEntity{
		Entity: types.Entity{
			ID:        id,
			CreatedAt: now,
		},
		CompactRequestStatistics: p,
	}
	return s.writeCompactStats(id, e)
}
func (s *BBolter) writeCompactStats(id string, e types.StatEntity) error {
	err := s.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket(BucketStats)
		bytes, err := s.Marshal(e)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(e.ID), bytes)
	})
	fmt.Println("updating stat", err)
	if err != nil {
		s.l.Error().Err(err).Msg("failed to save stats")
	} else {
		if e.UpdatedAt != nil {
			s.PublishChange(PubTypeStat, PubVerbUpdate, e)
		} else {
			s.PublishChange(PubTypeStat, PubVerbCreate, e)
		}
	}

	return err
}
func (s *BBolter) CleanCompactStats() error {
	err := s.emptyBucket(BucketStats)
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to empty bucket")
		return err
	}
	s.PublishChange(PubTypeStat, PubVerbClean, nil)
	// This is not really required, but can be helpful.
	s.compactDatabase()
	return err
}

func (s *BBolter) CompactStats() (es map[string]types.StatEntity, err error) {
	es = map[string]types.StatEntity{}
	err = s.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket(BucketStats)

		c := b.Cursor()

		for k, v := c.First(); k != nil; k, v = c.Next() {
			var e types.StatEntity
			err := s.Unmarshal(v, &e)
			if err != nil {
				return err
			}
			es[string(k)] = e
		}

		return nil
	})
	if err != nil {
		s.l.Error().Err(err).Msg("failed to lookup stats")
	}

	return es, err
}
