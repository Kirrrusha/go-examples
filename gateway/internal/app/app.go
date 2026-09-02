package app

import (
	"context"
	"fmt"
	"net"

	"gateway/internal/account"
	"gateway/internal/auth"
	"gateway/internal/config"
	"gateway/internal/interceptor"
	"gateway/internal/server"
	"gateway/internal/service"
	accountpb "gateway/pkg/account/go"
	authpb "gateway/pkg/auth/go"
	gatewaypb "gateway/pkg/gateway/go"
	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	cfg            *config.Config
	logger         *zerolog.Logger
	service        *service.GatewayService
	authService    *auth.Service
	accountService *account.Service
	server         *server.Server
	grpcServer     *grpc.Server
}

func New(logger *zerolog.Logger, cfg *config.Config) *App {
	return &App{
		cfg:    cfg,
		logger: logger,
	}
}

func (a *App) Run(ctx context.Context) error {
	gatewayServer, err := a.getGatewayServer(ctx)
	if err != nil {
		return fmt.Errorf("failed to get gateway server: %w", err)
	}

	a.grpcServer = a.getGRPCServer(gatewayServer)

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

func (a *App) getGatewayService(ctx context.Context) (*service.GatewayService, error) {
	if a.service == nil {
		authService, err := a.getAuthService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get auth service: %w", err)
		}

		accountService, err := a.getAccountService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get account service: %w", err)
		}

		a.service = service.New(accountService, authService, a.logger)
	}

	return a.service, nil
}

func (a *App) getAuthService(ctx context.Context) (*auth.Service, error) {
	if a.authService == nil {
		conn, err := grpc.NewClient(
			a.cfg.AuthGrpcHost,
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create auth service: %w", err)
		}

		client := authpb.NewAuthClient(conn)
		a.authService = auth.NewService(client)
	}

	return a.authService, nil
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

func (a *App) getGatewayServer(ctx context.Context) (*server.Server, error) {
	if a.server == nil {
		service, err := a.getGatewayService(ctx)
		if err != nil {
			return nil, fmt.Errorf("failed to get service: %w", err)
		}

		a.server = server.New(service, a.logger)
	}

	return a.server, nil
}

func (a *App) getGRPCServer(srv *server.Server) *grpc.Server {
	jwtInterceptor := interceptor.NewJWTInterceptor(a.cfg.JwtSecret)
	grpcSrv := grpc.NewServer(
		grpc.UnaryInterceptor(jwtInterceptor.UnaryInterceptor()),
		grpc.StreamInterceptor(jwtInterceptor.StreamInterceptor()),
	)
	gatewaypb.RegisterGatewayServer(grpcSrv, srv)

	return grpcSrv
}

func (a *App) Close() error {
	if a.grpcServer != nil {
		a.grpcServer.GracefulStop()
	}

	return nil
}
