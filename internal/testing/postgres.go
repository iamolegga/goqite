package testing

import (
	"context"
	"database/sql"
	"math/rand/v2"
	"strings"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

// NewPostgreSQLDB for testing.
func NewPostgreSQLDB(t *testing.T) *sql.DB {
	t.Helper()

	adminDB, adminClose := connect(t, "postgres")

	name := createName(t)
	if _, err := adminDB.ExecContext(t.Context(), `create database `+name); err != nil {
		t.Fatal(err)
	}
	db, close := connect(t, name)

	t.Cleanup(func() {
		close(t)
		if _, err := adminDB.ExecContext(context.WithoutCancel(t.Context()), `drop database if exists `+name); err != nil {
			t.Fatal(err)
		}
		adminClose(t)
	})

	return db
}

func connect(t *testing.T, name string) (*sql.DB, func(t *testing.T)) {
	t.Helper()

	db, err := sql.Open("pgx", "postgres://test:test@localhost:5433/"+name)
	if err != nil {
		t.Fatal(err)
	}

	return db, func(t *testing.T) {
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}

func createName(t *testing.T) string {
	t.Helper()

	const letters = "abcdefghijklmnopqrstuvwxyz"
	var b strings.Builder
	for range 16 {
		i := rand.IntN(len(letters))
		b.WriteByte(letters[i])
	}

	return b.String()
}
