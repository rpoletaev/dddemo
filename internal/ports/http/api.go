package http

import (
	"net/http"
	"path"
	"subscription/gen/http/restapi"
	"subscription/gen/http/restapi/operations"
	"subscription/internal/app/commands"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime/middleware"
	"github.com/rs/zerolog"
)

type Config struct {
	Port int `envconfig:"port"`
}

type Api struct {
	logger     zerolog.Logger
	config     *Config
	uowFactory commands.UnitOfWorkFactory
	server     *restapi.Server `wire:"-"`
}

func (api *Api) Connect() error {
	spec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		return err
	}

	subApi := operations.NewSubscriptionsServiceAPI(spec)
	server := restapi.NewServer(subApi)
	server.Port = api.config.Port

	api.server = server

	subApi.CreateSubscriptionHandler = operations.CreateSubscriptionHandlerFunc(api.CreateSubscription)

	var title string
	if sp := spec.Spec(); sp != nil && sp.Info != nil && sp.Info.Title != "" {
		title = sp.Info.Title
	}
	subApi.Middleware = func(b middleware.Builder) http.Handler {
		redocOpts := middleware.RedocOpts{
			BasePath: spec.BasePath(),
			SpecURL:  path.Join(spec.BasePath(), "/swagger.json"),
			Title:    title,
		}
		return middleware.Spec(spec.BasePath(), spec.Raw(), middleware.Redoc(redocOpts, subApi.Context().RoutesHandler(b)))
	}

	server.SetHandler(api.setupGlobalMiddleware(subApi.Serve(nil)))
	return nil
}

func (api *Api) Serve() error {
	return api.server.Serve()
}

func (api *Api) Shutdown() error {
	return api.server.Shutdown()
}

func (api *Api) Handler() http.Handler {
	return api.server.GetHandler()
}

func (api *Api) Close() {
	if err := api.server.Shutdown(); err != nil {
		api.logger.Fatal().Err(err)
	}
}
