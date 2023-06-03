package db

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type TestDB struct {
	m *migrate.Migrate
}

func NewTestDBMigrator() (*TestDB, error) {
	if os.Getenv("DB_URL") == "" {
		return nil, errors.New("needs TEST_DB_URL set on env")
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	prodapiIndex := strings.Index(cwd, "prodapi")
	if prodapiIndex == -1 {
		fmt.Println("prodapi not found in string")
		return nil, errors.New("incorrect project naming")
	}
	base := cwd[:prodapiIndex+len("prodapi")]

	// Construct the path to your migrations directory
	migrationsDir := filepath.Join(base, "db", "migrations")

	fmt.Println(migrationsDir)

	m, err := migrate.New(
		fmt.Sprintf("file://%s", migrationsDir),
		os.Getenv("DB_URL"))
	if err != nil {
		return nil, err
	}

	return &TestDB{
		m: m,
	}, nil
}

func (t TestDB) Prepare() error {
	err := t.m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		} else {
			return err
		}
	}
	return nil
}

func (t TestDB) Cleanup() error {
	err := t.m.Drop()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		} else {
			return err
		}
	}
	return nil
}
