package drivers

import (
	"context"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PgxDriver struct {
	dbURL string
	conn  *pgxpool.Pool
}

// Декораторы, чтобы логировать SQL
func (d *PgxDriver) exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	slog.Debug(sql, slog.Any("args", args))
	return d.conn.Exec(ctx, sql, args...)
}

func (d *PgxDriver) queryRows(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	slog.Debug(sql, slog.Any("args", args))
	return d.conn.Query(ctx, sql, args...)
}

func (d *PgxDriver) queryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	slog.Debug(sql, slog.Any("args", args))
	return d.conn.QueryRow(ctx, sql, args...)
}

func (d *PgxDriver) Open() error {
	pool, err := pgxpool.New(context.Background(), d.dbURL)
	if err != nil {
		return err
	}
	d.conn = pool

	var errConn error
	var ok bool
	for i := 1; i <= 5; i += 2 {
		if errConn = d.Ping(); errConn == nil {
			ok = true
			break
		}
		slog.Debug("try reconnect to database", slog.Int("sleep", i))
		time.Sleep(time.Duration(i) * time.Second)
	}
	if !ok {
		return errConn
	}

	return nil
}

func (d *PgxDriver) Close() error {
	d.conn.Close()
	return nil
}

func (d *PgxDriver) Ping() error {
	return d.conn.Ping(context.Background())
}
