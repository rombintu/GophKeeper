package drivers

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	apb "github.com/rombintu/GophKeeper/internal/proto/auth"
	kpb "github.com/rombintu/GophKeeper/internal/proto/keeper"
)

type key string

const (
	testKey key = "test"
)

type PgxDriver struct {
	dbURL       string
	serviceName string
	conn        *pgxpool.Pool
}

func NewPgxDriver(dbURL, serviceName string) *PgxDriver {
	return &PgxDriver{dbURL: dbURL, serviceName: serviceName}
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

func (d *PgxDriver) Open(ctx context.Context) error {
	pool, err := pgxpool.New(ctx, d.dbURL)
	if err != nil {
		return err
	}
	d.conn = pool

	var errConn error
	var ok bool
	for i := 1; i <= 5; i += 2 {
		if errConn = d.Ping(ctx, false); errConn == nil {
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

func (d *PgxDriver) Close(ctx context.Context) error {
	d.conn.Close()
	return nil
}

func (d *PgxDriver) Ping(ctx context.Context, monitoring bool) error {
	if err := d.conn.Ping(ctx); err != nil {
		return err
	}
	if monitoring {
		sql := `INSERT INTO services (name, last_check) VALUES ($1, NOW())
		ON CONFLICT (name) DO UPDATE SET last_check=NOW()`
		if _, err := d.exec(ctx, sql, d.serviceName); err != nil {
			return err
		}
	}
	return nil
}

func (d *PgxDriver) UserGet(ctx context.Context, user *apb.User) (*apb.User, error) {
	sql := `SELECT email, key_checksum FROM users WHERE email=$1`
	row := d.queryRow(ctx, sql, user.GetEmail())
	var findUser apb.User

	if err := row.Scan(&findUser.Email, &findUser.KeyChecksum); err != nil {
		return nil, err
	}
	return &findUser, nil
}

func (d *PgxDriver) UserCreate(ctx context.Context, user *apb.User) error {
	sql := `INSERT INTO users (email, key_checksum) VALUES ($1, $2)`
	if _, err := d.exec(ctx, sql, user.GetEmail(), user.GetKeyChecksum()); err != nil {
		return err
	}
	return nil
}

// Переделать все под ID для надежности
func (d *PgxDriver) SecretCreate(ctx context.Context, secret *kpb.Secret) error {
	sql := `INSERT INTO secrets (
		title, secret_type, user_email, version, hash_payload, payload
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (title) DO UPDATE SET version = secrets.version+1`
	if _, err := d.exec(ctx, sql,
		secret.GetTitle(), secret.GetSecretType(),
		secret.GetUserEmail(), secret.GetVersion(),
		secret.GetPayload(),
	); err != nil {
		return err
	}
	return nil
}

func (d *PgxDriver) SecretCreateBatch(ctx context.Context, secrets []*kpb.Secret) error {
	tx, err := d.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback(ctx) //nolint:errcheck
	}()

	sql := `INSERT INTO secrets (
		title, secret_type, user_email, version, hash_payload, payload
		) VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (title) DO UPDATE SET version = secrets.version+1`

	var errs []error
	for _, s := range secrets {
		if _, err := d.exec(ctx, sql,
			s.GetTitle(), s.GetSecretType(),
			s.GetUserEmail(), s.GetVersion(),
			s.GetHashPayload(), s.GetPayload(),
		); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

func (d *PgxDriver) SecretList(ctx context.Context, userEmail string) ([]*kpb.Secret, error) {
	sql := `SELECT title, secret_type, user_email, created_at, version, hash_payload, payload 
		FROM secrets s 
		INNER JOIN (
		SELECT id, MAX(version) AS max_version
			FROM secrets
			GROUP BY id
		) AS sub
		ON s.id = sub.id AND s.version = sub.max_version
		WHERE user_email=$1
		`
	rows, err := d.queryRows(ctx, sql, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var secrets []*kpb.Secret
	for rows.Next() {
		var createdAt time.Time
		var s kpb.Secret
		if err := rows.Scan(
			&s.Title, &s.SecretType, &s.UserEmail, &createdAt, &s.Version, &s.HashPayload, &s.Payload,
		); err != nil {
			return nil, err
		}
		// Преобразование time.Time в google.protobuf.Timestamp
		pbTimestamp := &timestamp.Timestamp{
			Seconds: createdAt.Unix(),
			Nanos:   int32(createdAt.Nanosecond()),
		}
		s.CreatedAt = pbTimestamp
		secrets = append(secrets, &s)
	}
	return secrets, nil
}

func (d *PgxDriver) SecretPurge(ctx context.Context, secret *kpb.Secret) error {
	sql := `DELETE FROM secrets WHERE title=$1, user_email=$2`
	if _, err := d.exec(ctx, sql,
		secret.GetTitle(), secret.GetUserEmail(),
	); err != nil {
		return err
	}
	return nil
}

func (d *PgxDriver) autoDefaultMigrate(mpath string) error {

	migr, err := migrate.New(
		fmt.Sprintf("file://%s", mpath),
		d.dbURL,
	)
	if err != nil {
		return err
	}
	return migr.Up()
}

func (d *PgxDriver) Configure(ctx context.Context) error {
	mpath, err := filepath.Abs(
		filepath.Join("internal", "storage", "migrations"))
	if err != nil {
		return err
	}

	slog.Debug("migration init", slog.String("path", mpath))
	istest := ctx.Value(testKey)
	if istest != nil && istest == true {
		// Получаем абсолютный путь к директории migrations
		mpath, err = filepath.Abs(filepath.Join("..", "migrations"))
		if err != nil {
			return fmt.Errorf("failed to get absolute path for migrations: %w", err)
		}
	}
	if err := d.autoDefaultMigrate(mpath); err != nil {
		slog.Warn("auto migration failed or skip", slog.String("message", err.Error()))
	}
	return nil
}

func (d *PgxDriver) SecretGetBatch(ctx context.Context) ([]*kpb.Secret, error) {
	return nil, nil
}
