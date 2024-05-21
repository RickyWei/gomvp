package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"

	"github.com/rickywei/sparrow/project/api/middleware"
	"github.com/rickywei/sparrow/project/conf"
	"github.com/rickywei/sparrow/project/graph"
	"github.com/rickywei/sparrow/project/graph/resolver"
)

var (
	ProviderSet = wire.NewSet(NewApi)
)

type API struct {
	engine *gin.Engine
	srv    *http.Server
	ctx    context.Context
	cancel context.CancelFunc
}

func NewApi() *API {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)

	engine := gin.New()
	engine.Use(middleware.Logger(),
		middleware.Recover(),
		middleware.GinContextToContext(),
	)

	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", conf.String("app.ip"), conf.Int("app.port")),
		Handler: engine,
	}

	api := &API{
		engine: engine,
		srv:    srv,
		ctx:    ctx,
		cancel: cancel,
	}
	api.engine.POST("/graphql", graphqlHandler())
	if !conf.IsProd() {
		api.engine.GET("/playground", playgroundHandler())
	}

	return api
}

func (a *API) Run() (err error) {
	return a.srv.ListenAndServe()
}

func (a *API) Stop() {
	defer a.cancel()

	a.srv.Shutdown(a.ctx)
}

func graphqlHandler() gin.HandlerFunc {
	h := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &resolver.Resolver{}}))
	h.Use(extension.AutomaticPersistedQuery{Cache: middleware.GetApqCache()})

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func playgroundHandler() gin.HandlerFunc {
	h := playground.Handler("GraphQL", "/graphql")

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
