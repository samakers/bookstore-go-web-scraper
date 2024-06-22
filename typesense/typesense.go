package typesensehandler

import (
	"context"

	"github.com/typesense/typesense-go/typesense"
	"github.com/typesense/typesense-go/typesense/api"
	"github.com/typesense/typesense-go/typesense/api/pointer"
)

func typesensehandler() {
	client := typesense.NewClient(
		typesense.WithServer("http://localhost:8108"),
		typesense.WithAPIKey("Hu52dwsas2AdxdE"))

	schema := &api.CollectionSchema{
		Name: "books",
		Fields: []api.Field{
			{
				Name: "title",
				Type: "string",
			},
			{
				Name: "price",
				Type: "string",
			},
			{
				Name: "availability",
				Type: "string",
			},
		},
		DefaultSortingField: pointer.String("title"),
	}

	client.Collections().Create(context.Background(), schema)

}
