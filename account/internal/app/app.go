package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	"account/internal/config"
	"account/internal/repository"
	"account/internal/server"
	"account/internal/service"
	accountpb "account/pkg/account/go"
	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type App struct {
	cfg               *config.Config
	logger            *zerolog.Logger
	accountRepository *repository.Repository
	accountService    *service.AccountService
	accountServer     *server.Server
	grpcServer        *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	accountServer, err := a.getAccountServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get account server: %w", err)
	}

	a.grpcServer = getGRPCServer(accountServer)

	listenAddr := fmt.Sprintf("%s:%d", a.cfg.Host, a.cfg.Port)
	lis, err := net.Listen("tcp", listenAddr)
	if err != nil {
		a.logger.Fatal().Err(err).Msg("failed to listen")
		return err
	}

	a.logger.Info().Str("addr", listenAddr).Msg("gRPC server listening")

	serveErrCh := make(chan error, 1)
	go func() {
		serveErrCh <- a.grpcServer.Serve(lis)
	}()

	select {
	case <-ctx.Done():
		a.grpcServer.GracefulStop()
		return ctx.Err()
	case err := <-serveErrCh:
		if err != nil {
			a.logger.Error().Err(err).Msg("failed to serve")
			return err
		}
	}

	return nil
}

func (a *App) getRepository(ctx context.Context) (*repository.Repository, error) {
	if a.accountRepository == nil {
		if err := a.runMigrations(ctx); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}

		db, err := gorm.Open(postgres.Open(a.cfg.DbDsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("gorm init failed: %w", err)
		}

		a.accountRepository = repository.NewRepository(db, a.logger)
	}

	return a.accountRepository, nil
}

func (a *App) runMigrations(ctx context.Context) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("failed to select migrations dialect: %w", err)
	}

	dbGoose, err := sql.Open("postgres", a.cfg.DbDsn)
	if err != nil {
		return fmt.Errorf("failed to create sql connection: %w", err)
	}
	defer dbGoose.Close()

	if err := goose.UpContext(ctx, dbGoose, "migrations"); err != nil {
		return fmt.Errorf("failed to run up migrations: %w", err)
	}

	return nil
}

func (a *App) getAccountService(ctx context.Context) (*service.AccountService, error) {
	if a.accountService == nil {
		repo, err := a.getRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}

		a.accountService = service.New(repo, a.logger)
	}

	return a.accountService, nil
}

func (a *App) getAccountServer(ctx context.Context) (*server.Server, error) {
	if a.accountServer == nil {
		service, err := a.getAccountService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get service: %w", err)
		}

		a.accountServer = server.New(service, a.logger)
	}

	return a.accountServer, nil
}

func getGRPCServer(srv *server.Server) *grpc.Server {
	grpcSrv := grpc.NewServer()
	accountpb.RegisterAccountServer(grpcSrv, srv)
	return grpcSrv
}

func (a *App) Close() error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	return nil
}
