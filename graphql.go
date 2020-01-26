package main

import (
	"fmt"

	"github.com/graphql-go/graphql"
	"github.com/minheq/kedulv2/service_salon/account"
	"github.com/minheq/kedulv2/service_salon/business"
)

type gql struct {
	accountService  account.Service
	businessService business.Service
}

func (g *gql) schema() graphql.Schema {
	businessType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Business",
			Fields: graphql.Fields{
				"id": &graphql.Field{
					Type: graphql.Int,
				},
			},
		},
	)

	queryType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"business": &graphql.Field{
					Type: businessType,
					Args: graphql.FieldConfigArgument{
						"id": &graphql.ArgumentConfig{
							Type: graphql.Int,
						},
					},

					Resolve: func(p graphql.ResolveParams) (interface{}, error) {
						id, ok := p.Args["id"].(string)

						if ok != true {
							return nil, fmt.Errorf("sux")
						}

						business, err := g.businessService.GetByID(id)

						if err != nil {
							return nil, err
						}

						return business, err
					},
				},
			},
		})

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"create": &graphql.Field{
				Type: businessType,
				Args: graphql.FieldConfigArgument{
					"name": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.String),
					},
				},
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					business, err := g.businessService.CreateBusiness(business.CreateBusinessInput{
						Name: fmt.Sprint(params.Args["name"]),
					})

					if err != nil {
						return nil, err
					}

					return business, nil
				},
			},
		},
	})

	schema, _ := graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)

	return schema
}
