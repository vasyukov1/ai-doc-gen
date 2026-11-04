package storage

import (
	"context"
	"database/sql/driver"
	"errors"
	"time"

	"github.com/go-jet/jet/v2/postgres"
)

// A session lock prevents cleanup while any process still uses job directories.
func (r *Repository) GuardJob(ctx context.Context, id string) (func(), error) {
	acquire, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	conn, err := r.db.Conn(acquire)
	if err != nil {
		return nil, err
	}
	query, args := postgres.SELECT(postgres.RawBool("pg_try_advisory_lock(hashtextextended(:id, 0))", postgres.RawArgs{":id": id})).Sql()
	var locked bool
	if err = conn.QueryRowContext(acquire, query, args...).Scan(&locked); err != nil || !locked {
		if err != nil {
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		}
		conn.Close()
		if err == nil {
			err = errors.New("job process still active")
		}
		return nil, err
	}
	return func() {
		release, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		query, args := postgres.SELECT(postgres.RawBool("pg_advisory_unlock(hashtextextended(:id, 0))", postgres.RawArgs{":id": id})).Sql()
		var unlocked bool
		if err := conn.QueryRowContext(release, query, args...).Scan(&unlocked); err != nil || !unlocked {
			_ = conn.Raw(func(any) error { return driver.ErrBadConn })
		}
		_ = conn.Close()
	}, nil
}
