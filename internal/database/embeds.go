package database

import (
	"embed"
	"io/fs"
)

//go:embed migrations
var migrationsFS embed.FS

func MigrationFS() fs.FS {
	sub, err := fs.Sub(migrationsFS, "migrations")
	if err != nil {
		panic(err)
	}
	return sub
}
