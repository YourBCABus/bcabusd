package api

import (
	"github.com/go-pg/pg/v9"
	"github.com/graphql-go/graphql"
)

func makeQuery(db *pg.DB) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{
		"apiVersion": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "0.1", nil
			},
		},
		"lifeTheUniverseAndEverything": &graphql.Field{
			Type: graphql.Int,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				var answer int
				_, err := db.Query(pg.Scan(&answer), "SELECT 42")
				if err != nil {
					return nil, err
				}
				return answer, nil
			},
		},
	}})
}

// MakeSchema creates the YourBCABus API
// GraphQL schema, given a Postgres database.
func MakeSchema(db *pg.DB) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{Query: makeQuery(db)})
}
