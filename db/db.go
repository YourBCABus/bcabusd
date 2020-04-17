package db

import (
	"context"
	"fmt"
	"os"

	"github.com/go-pg/pg/v9"
)

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	fmt.Println(q.FormattedQuery())
	return nil
}

// Connect connects to the database specified in
// the DATABASE_URL environment variable.
func Connect() (*pg.DB, error) {
	options, err := pg.ParseURL(os.Getenv("DATABASE_URL"))

	if err != nil {
		return nil, err
	}

	db := pg.Connect(options)
	db.AddQueryHook(dbLogger{})
	return db, nil
}
