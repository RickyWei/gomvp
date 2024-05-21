//go:build wireinject

package app

import (
	"github.com/google/wire"

	"github.com/rickywei/sparrow/project/api"
)

func WireApp() (*App, error) {
	panic(wire.Build(api.ProviderSet, NewApp))
}
