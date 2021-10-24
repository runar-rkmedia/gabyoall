package bboltStorage

import (
	"fmt"
	"os"
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
	if err != nil {
		s.l.Error().Err(err).Msg("failed to save stats")
	}
	return err
}
func (s *BBolter) CleanCompactStats() error {
	err := s.emptyBucket(BucketStats)
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to empty bucket")
		return err
	}
	// This is not really required, but can be helpful.
	s.compactDatabase()
	return err
}
func (s *BBolter) compactDatabase() error {
	path := "_compact.bbolt"
	originalPath := s.Path()
	compactDb, err := bolt.Open(path, 0666, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to create a new database before compacting this one")
		return fmt.Errorf("Failed to create a new database before compacting this one")
	}
	s.l.Warn().Msg("New database was opened")
	err = bolt.Compact(compactDb, s.DB, 0)
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to compact database")
		return fmt.Errorf("Failed to compact database")
	}
	s.l.Warn().Msg("New database was compacted. Will now close existing database.")
	compactDb.Close()
	s.Close()
	s.l.Warn().Msg("Closed databases. Will now rename databases on disk")
	err = os.Rename(originalPath, originalPath+".bk")
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to move original database")
		return fmt.Errorf("Failed to move original database")
	}
	err = os.Rename(path, originalPath)
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to move compact database")
		return fmt.Errorf("Failed to move compact database")
	}
	s.l.Warn().Msg("Databases renamed. WIll now reopen the database.")
	db, err := bolt.Open(originalPath, 0666, &bolt.Options{
		Timeout: 1 * time.Second,
	})
	if err != nil {
		s.l.Error().Err(err).Msg("Failed to reopen the database")
		return fmt.Errorf("Failed to reopen the database")
	}
	s.DB = db

	s.l.Info().Msg("Database was compacted and replaced successfully")
	return nil
}
func (s *BBolter) emptyBucket(bucket []byte) error {
	s.l.Warn().Str("bucket", string(bucket)).Msg("Emptying bucket")
	err := s.Update(func(tx *bolt.Tx) error {
		err := tx.DeleteBucket(bucket)
		if err != nil {
			return fmt.Errorf("Failed to delete bucket %s: %w", string(bucket), err)
		}
		_, err = tx.CreateBucketIfNotExists(bucket)
		if err != nil {
			return fmt.Errorf("Failed to reacreate bucket %s: %w", string(bucket), err)
		}
		return nil
	})
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
