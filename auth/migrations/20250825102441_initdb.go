package migrations

import (
	"context"
	"database/sql"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddNamedMigrationContext("20250825102441_initdb.go", upInitDB, downInitDB)
}

func upInitDB(ctx context.Context, tx *sql.Tx) error {
	query := `
CREATE TABLE IF NOT EXISTS users (
	id BIGSERIAL PRIMARY KEY,
	login TEXT NOT NULL UNIQUE,
	email TEXT NOT NULL UNIQUE,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW(),
	updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	query = `
CREATE TABLE IF NOT EXISTS refresh_tokens (
	id BIGSERIAL PRIMARY KEY,
	user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
	token TEXT NOT NULL UNIQUE,
	expires_at TIMESTAMP NOT NULL,
	revoked_at TIMESTAMP NULL,
	created_at TIMESTAMP NOT NULL DEFAULT NOW()
);`

	if _, err := tx.ExecContext(ctx, query); err != nil {
		return err
	}

	return nil
}

func downInitDB(ctx context.Context, tx *sql.Tx) error {
	if _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS refresh_tokens;`); err != nil {
		return err
	}

	if _, err := tx.ExecContext(ctx, `DROP TABLE IF EXISTS users;`); err != nil {
		return err
	}

	return nil
}
