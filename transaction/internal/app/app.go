package app

import (
	"context"
	"database/sql"
	"fmt"
	"net"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"transaction/internal/account"
	"transaction/internal/config"
	"transaction/internal/repository"
	"transaction/internal/server"
	"transaction/internal/service"
	_ "transaction/migrations"
	accountpb "transaction/pkg/account/go"
	transactionpb "transaction/pkg/transaction/go"
)

type App struct {
	cfg                   *config.Config
	logger                *zerolog.Logger
	transactionRepository *repository.Repository
	accountService        *account.Service
	transactionService    *service.TransactionService
	transactionServer     *server.Server
	grpcServer            *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{cfg: cfg, logger: logger}
}

func (a *App) Run(ctx context.Context) error {
	transactionServer, err := a.getTransactionServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get transaction server: %w", err)
	}

	a.grpcServer = getGRPCServer(transactionServer)

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
	if a.transactionRepository == nil {
		if err := a.runMigrations(ctx); err != nil {
			return nil, fmt.Errorf("failed to run migrations: %w", err)
		}

		db, err := gorm.Open(postgres.Open(a.cfg.DbDsn), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("gorm init failed: %w", err)
		}

		a.transactionRepository = repository.NewRepository(db, a.logger)
	}

	return a.transactionRepository, nil
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

func (a *App) getAccountService(ctx context.Context) (*account.Service, error) {
	if a.accountService == nil {
		conn, err := grpc.NewClient(
			a.cfg.AccountGrpcHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create account service: %w", err)
		}

		client := accountpb.NewAccountClient(conn)
		a.accountService = account.New(client)
	}

	return a.accountService, nil
}

func (a *App) getTransactionService(ctx context.Context) (*service.TransactionService, error) {
	if a.transactionService == nil {
		repo, err := a.getRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get repository: %w", err)
		}

		accountService, err := a.getAccountService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get account service: %w", err)
		}

		a.transactionService = service.New(repo, accountService, a.logger)
	}

	return a.transactionService, nil
}

func (a *App) getTransactionServer(ctx context.Context) (*server.Server, error) {
	if a.transactionServer == nil {
		service, err := a.getTransactionService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get service: %w", err)
		}

		a.transactionServer = server.New(service, a.logger)
	}

	return a.transactionServer, nil
}

func getGRPCServer(srv *server.Server) *grpc.Server {
	grpcSrv := grpc.NewServer()
	transactionpb.RegisterTransactionServiceServer(grpcSrv, srv)
	return grpcSrv
}

func (a *App) Close() error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	return nil
}
