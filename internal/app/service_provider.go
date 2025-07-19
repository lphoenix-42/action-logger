package app

import (
	"context"
	"log"
	"net/http"

	"github.com/lphoenix-42/action-logger/gen/actionlog/v1/actionlog_v1connect"
	"github.com/lphoenix-42/action-logger/internal/config"
	"github.com/lphoenix-42/action-logger/internal/config/env"
	actionlogAPI "github.com/lphoenix-42/action-logger/internal/infrastructure/delivery/actionlog"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type serviceProvider struct {
	httpConfig config.HTTPConfig

	actionlogAPI *actionlogAPI.Server
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) ActionlogAPI(ctx context.Context) *actionlogAPI.Server {
	if s.actionlogAPI == nil {
		s.actionlogAPI = actionlogAPI.New()
	}

	return s.actionlogAPI
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()

	path, handler := actionlog_v1connect.NewActionLogServiceHandler(a.serviceProvider.ActionlogAPI(ctx))
	mux.Handle(path, handler)

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}
	return nil

}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func (s *serviceProvider) HTTPConfig() config.HTTPConfig {
	if s.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to get http config: %s", err.Error())
		}

		s.httpConfig = cfg
	}

	return s.httpConfig
}
