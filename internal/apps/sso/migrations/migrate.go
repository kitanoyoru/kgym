package migrations

import (
	"context"
	"path/filepath"
	"runtime"

	"github.com/pressly/goose/v3"
)

func getMigrationsDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Dir(filename)
}

func Up(ctx context.Context, driver, uri string) error {
	migrationsDir := getMigrationsDir()

	sql, err := goose.OpenDBWithDriver(driver, uri)
	if err != nil {
		return err
	}
	defer sql.Close()

	err = goose.Up(sql, migrationsDir)
	if err != nil {
		return err
	}

	return nil
}

func Down(ctx context.Context, driver, uri string) error {
	migrationsDir := getMigrationsDir()

	sql, err := goose.OpenDBWithDriver(driver, uri)
	if err != nil {
		return err
	}
	defer sql.Close()

	err = goose.Down(sql, migrationsDir)
	if err != nil {
		return err
	}

	return nil
}
