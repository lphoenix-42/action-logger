package app

import (
	"context"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lphoenix-42/action-logger/gen/actionlog/v1/actionlog_v1connect"
	"github.com/lphoenix-42/action-logger/internal/config"
	"github.com/lphoenix-42/action-logger/internal/config/env"
	actionlogAPI "github.com/lphoenix-42/action-logger/internal/infrastructure/delivery/actionlog"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/repository"
	actionlogRepositoryPg "github.com/lphoenix-42/action-logger/internal/infrastructure/repository/actionlog/pg"
	"github.com/lphoenix-42/action-logger/internal/service"
	actionlogService "github.com/lphoenix-42/action-logger/internal/service/actionlog"
	"github.com/lphoenix-42/action-logger/pkg/closer"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type serviceProvider struct {
	httpConfig config.HTTPConfig
	pgConfig   config.PGConfig

	actionlogAPI        *actionlogAPI.Server
	actionlogService    service.ActionlogService
	actionlogRepository repository.ActionlogRepository

	pgConn *pgxpool.Pool
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (s *serviceProvider) ActionlogAPI(ctx context.Context) *actionlogAPI.Server {
	if s.actionlogAPI == nil {
		s.actionlogAPI = actionlogAPI.New(s.ActionlogService(ctx))
	}

	return s.actionlogAPI
}

func (s *serviceProvider) ActionlogService(ctx context.Context) service.ActionlogService {
	if s.actionlogService == nil {
		s.actionlogService = actionlogService.New(s.ActionlogRepository(ctx))
	}

	return s.actionlogService
}

func (s *serviceProvider) ActionlogRepository(ctx context.Context) repository.ActionlogRepository {
	if s.actionlogRepository == nil {
		s.actionlogRepository = actionlogRepositoryPg.New(s.PGConn(ctx))
	}

	return s.actionlogRepository
}

func (s *serviceProvider) PGConn(ctx context.Context) *pgxpool.Pool {
	if s.pgConn == nil {
		pool, err := pgxpool.New(ctx, s.PGConfig().DSN())
		if err != nil {
			log.Fatalf("failed to connect to database: %v", err)
		}

		err = pool.Ping(ctx)
		if err != nil {
			log.Fatalf("ping error: %s", err.Error())
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		s.pgConn = pool
	}

	return s.pgConn

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

func (s *serviceProvider) PGConfig() config.PGConfig {
	if s.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to get pg config: %s", err.Error())
		}

		s.pgConfig = cfg
	}

	return s.pgConfig
}
