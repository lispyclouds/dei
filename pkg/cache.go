package pkg

import (
	"context"
	"os"
	"path"

	"database/sql"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type Cache struct {
	db *sql.DB
}

type WriteTxn struct {
	cache *Cache
	pairs []string
}

func NewCache() (*Cache, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return nil, err
	}

	dbDir := path.Join(cacheDir, "dei")
	dbPath := path.Join(dbDir, "cache.sqlite.db")

	if _, err := os.Stat(dbPath); !os.IsNotExist(err) {
		db, err := sql.Open("sqlite3", dbPath)
		if err != nil {
			return nil, err
		}

		return &Cache{db: db}, nil
	}

	if err = os.MkdirAll(dbDir, os.ModePerm); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS cache (key TEXT PRIMARY KEY, value TEXT)")
	if err != nil {
		return nil, err
	}

	return &Cache{db: db}, nil
}

func insert(db *sql.DB, ctx context.Context, k, v string) error {
	_, err := db.ExecContext(ctx, "INSERT INTO cache VALUES (?, ?) ON CONFLICT (key) DO UPDATE SET value=?", k, v, v)
	return err
}

func (c *Cache) Close() error {
	return c.db.Close()
}

func (c *Cache) Put(k, v string) error {
	return insert(c.db, context.Background(), k, v)
}

func (c *Cache) Get(k string) (string, error) {
	result, err := c.db.Query("SELECT value FROM cache WHERE key=?", k)
	if err != nil {
		return "", err
	}

	if !result.Next() {
		return "", nil
	}

	var value string
	if err = result.Scan(&value); err != nil {
		return "", err
	}
	result.Close()

	return value, nil
}

func (c *Cache) WithWriteTxn() *WriteTxn {
	return &WriteTxn{cache: c}
}

func (wtx *WriteTxn) Put(k, v string) *WriteTxn {
	wtx.pairs = append(wtx.pairs, k, v)
	return wtx
}

func (wtx *WriteTxn) Run() error {
	ctx := context.Background()

	tx, err := wtx.cache.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for i := 0; i+1 < len(wtx.pairs); i += 2 {
		if err = insert(wtx.cache.db, ctx, wtx.pairs[i], wtx.pairs[i+1]); err != nil {
			return err
		}
	}

	return tx.Commit()
}
