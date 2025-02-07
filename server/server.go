package server

import (
	"context"
	"errors"
	"log"
	"myhub/graph"
	"myhub/internal/storage"
	"myhub/internal/utils"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const defaultPort = "8080"

func StartServer(dataStore storage.DataStore) {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{
		Resolvers: &graph.Resolver{Storage: dataStore},
	}))
	srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
		formattedErr := graphql.DefaultErrorPresenter(ctx, err)

		var gqlErr *utils.GraphQLError
		if errors.As(err, &gqlErr) {
			formattedErr.Message = gqlErr.Message
			formattedErr.Extensions = map[string]interface{}{
				"code": gqlErr.Code,
			}
		} else {
			log.Printf("Неизвестная ошибка: %v", err)
			formattedErr.Message = "Internal server error"
			formattedErr.Extensions = map[string]interface{}{
				"code": "INTERNAL_ERROR",
			}
		}

		return formattedErr
	})
	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
