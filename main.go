package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/graphql-go/graphql"
)

// Product contains information about one product
type Product struct {
	ID    int64   `json:"id"`
	Name  string  `json:"name"`
	Info  string  `json:"info,omitempty"`
	Price float64 `json:"price"`
}

type Order struct {
	ID        int64     `json:"id"`
	OrderDate time.Time `json:"orderDate"`
	Products  []Product `json:"products"`
}

var products = []Product{
	{
		ID:    1,
		Name:  "Chicha Morada",
		Info:  "Chicha morada is a beverage originated in the Andean regions of Perú but is actually consumed at a national level (wiki)",
		Price: 7.99,
	},
	{
		ID:    2,
		Name:  "Chicha de jora",
		Info:  "Chicha de jora is a corn beer chicha prepared by germinating maize, extracting the malt sugars, boiling the wort, and fermenting it in large vessels (traditionally huge earthenware vats) for several days (wiki)",
		Price: 5.95,
	},
	{
		ID:    3,
		Name:  "Cisco",
		Info:  "Cisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
		Price: 34.95,
	},
	{
		ID:    4,
		Name:  "Xisco",
		Info:  "Xisco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
		Price: 23.95,
	},
	{
		ID:    5,
		Name:  "Disco",
		Info:  "Disco is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
		Price: 45.95,
	},
	{
		ID:    6,
		Name:  "Moook",
		Info:  "Moook is a colorless or yellowish-to-amber colored brandy produced in winemaking regions of Peru and Chile (wiki)",
		Price: 3.95,
	},
}

var orders = []Order{
	{
		ID:        1,
		OrderDate: time.Now(),
		Products:  []Product{products[0], products[2], products[3]},
	},
	{
		ID:        2,
		OrderDate: time.Now(),
		Products:  []Product{products[1], products[5]},
	},
	{
		ID:        3,
		OrderDate: time.Now(),
		Products:  []Product{products[1], products[3]},
	},
	{
		ID:        4,
		OrderDate: time.Now(),
		Products:  []Product{products[4], products[5]},
	},
	{
		ID:        5,
		OrderDate: time.Now(),
		Products:  []Product{products[2], products[0]},
	},
}

var productType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Product",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.Int,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"info": &graphql.Field{
				Type: graphql.String,
			},
			"price": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Product",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"orderdate": &graphql.Field{
			Type: graphql.DateTime,
		},
		"products": &graphql.Field{
			Type: graphql.NewList(productType),
		},
	},
})

var queryType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			/* Get (read) single product by id
			   http://localhost:8080/product?query={product(id:1){name,info,price}}
			*/
			"product": &graphql.Field{
				Type:        productType,
				Description: "Get product by id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						// Find product
						for _, product := range products {
							if int(product.ID) == id {
								return product, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) product list
			   http://localhost:8080/product?query={list{id,name,info,price}}
			*/
			"list": &graphql.Field{
				Type:        graphql.NewList(productType),
				Description: "Get product list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return products, nil
				},
			},
			/* Get (read) single order by id
			   http://localhost:8080/order?query={order(id:1){id,orderdate,products{id,price}}}
			*/
			"order": &graphql.Field{
				Type:        orderType,
				Description: "Get orders by Id",
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.Int,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["id"].(int)
					if ok {
						// Find product
						for _, order := range orders {
							if int(order.ID) == id {
								return order, nil
							}
						}
					}
					return nil, nil
				},
			},
			/* Get (read) product list
			   http://localhost:8080/order?query={orderlist{id,orderdate,products{id,price}}}
			*/
			"orderlist": &graphql.Field{
				Type:        graphql.NewList(orderType),
				Description: "Get order list",
				Resolve: func(params graphql.ResolveParams) (interface{}, error) {
					return orders, nil
				},
			},
		},
	})

var mutationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Mutation",
	Fields: graphql.Fields{
		/* Create new product item
		http://localhost:8080/product?query=mutation+_{create(name:"Inca Kola",info:"Inca Kola is a soft drink that was created in Peru in 1935 by British immigrant Joseph Robinson Lindley using lemon verbena (wiki)",price:1.99){id,name,info,price}}
		*/
		"create": &graphql.Field{
			Type:        productType,
			Description: "Create new product",
			Args: graphql.FieldConfigArgument{
				"name": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.String),
				},
				"info": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Float),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				rand.Seed(time.Now().UnixNano())
				product := Product{
					ID:    int64(rand.Intn(100000)), // generate random ID
					Name:  params.Args["name"].(string),
					Info:  params.Args["info"].(string),
					Price: params.Args["price"].(float64),
				}
				products = append(products, product)
				return product, nil
			},
		},

		/* Update product by id
		   http://localhost:8080/product?query=mutation+_{update(id:1,price:3.95){id,name,info,price}}
		*/
		"update": &graphql.Field{
			Type:        productType,
			Description: "Update product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
				"name": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"info": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
				"price": &graphql.ArgumentConfig{
					Type: graphql.Float,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				name, nameOk := params.Args["name"].(string)
				info, infoOk := params.Args["info"].(string)
				price, priceOk := params.Args["price"].(float64)
				product := Product{}
				for i, p := range products {
					if int64(id) == p.ID {
						if nameOk {
							products[i].Name = name
						}
						if infoOk {
							products[i].Info = info
						}
						if priceOk {
							products[i].Price = price
						}
						product = products[i]
						break
					}
				}
				return product, nil
			},
		},

		/* Delete product by id
		   http://localhost:8080/product?query=mutation+_{delete(id:1){id,name,info,price}}
		*/
		"delete": &graphql.Field{
			Type:        productType,
			Description: "Delete product by id",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.NewNonNull(graphql.Int),
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				id, _ := params.Args["id"].(int)
				product := Product{}
				for i, p := range products {
					if int64(id) == p.ID {
						product = products[i]
						// Remove from product list
						products = append(products[:i], products[i+1:]...)
					}
				}

				return product, nil
			},
		},
	},
})

var schema, _ = graphql.NewSchema(
	graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	},
)

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("errors: %v", result.Errors)
	}
	return result
}

func main() {
	http.HandleFunc("/product", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})
	http.HandleFunc("/order", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		json.NewEncoder(w).Encode(result)
	})

	fmt.Println("Server is running on port 8080")
	http.ListenAndServe(":8080", nil)
}
