package pkg

import (
	"os"
	"path"

	"go.etcd.io/bbolt"
)

type Cache struct {
	db     *bbolt.DB
	bucket []byte
}

type pair struct {
	k string
	v []byte
}

type WriteTxn struct {
	cache *Cache
	pairs []pair
}

func NewCache() (*Cache, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dbDir := path.Join(cacheDir, "dei")
	if err = os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := bbolt.Open(path.Join(dbDir, "cache.bolt.db"), 0600, nil)
	if err != nil {
		return nil, err
	}

	bucketName := []byte("cache")
	if err = db.Update(func(tx *bbolt.Tx) error {
		if bucket := tx.Bucket(bucketName); bucket == nil {
			_, err = tx.CreateBucket(bucketName)
			return err
		}

		return nil
	}); err != nil {
		return nil, err
	}

	return &Cache{db: db, bucket: bucketName}, nil
}

func (c *Cache) Close() error {
	return c.db.Close()
}

func (c *Cache) Put(k string, v []byte) error {
	return c.db.Update(func(tx *bbolt.Tx) error {
		return tx.Bucket(c.bucket).Put([]byte(k), v)
	})
}

func (c *Cache) Get(k string) ([]byte, error) {
	var value []byte

	if err := c.db.View(func(tx *bbolt.Tx) error {
		value = tx.Bucket(c.bucket).Get([]byte(k))
		return nil
	}); err != nil {
		return nil, err
	}

	return value, nil
}

func (c *Cache) WithWriteTxn() *WriteTxn {
	return &WriteTxn{cache: c}
}

func (wtx *WriteTxn) Put(k string, v []byte) *WriteTxn {
	wtx.pairs = append(wtx.pairs, pair{k, v})
	return wtx
}

func (wtx *WriteTxn) Run() error {
	return wtx.cache.db.Update(func(tx *bbolt.Tx) error {
		for _, pair := range wtx.pairs {
			if err := tx.Bucket(wtx.cache.bucket).Put([]byte(pair.k), pair.v); err != nil {
				return err
			}
		}

		return nil
	})
}
