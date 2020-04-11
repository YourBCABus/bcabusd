package api

import (
	"github.com/graphql-go/graphql"
)

func makeQuery() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{Name: "Query", Fields: graphql.Fields{
		"apiVersion": &graphql.Field{
			Type: graphql.String,
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				return "0.1", nil
			},
		},
	}})
}

func MakeSchema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{Query: makeQuery()})
}
