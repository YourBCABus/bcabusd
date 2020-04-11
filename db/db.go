package db

import (
	"os"

	"github.com/go-pg/pg/v9"
)

// Connect connects to the database specified in
// the DATABASE_URL environment variable.
func Connect() (*pg.DB, error) {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}

	return pg.Connect(options), nil
}
