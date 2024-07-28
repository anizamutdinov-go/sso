package grpcapp

import (
	"fmt"
	grpcauth "github.com/anizamutdinov-go/sso/internal/grpc/auth"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gPRCServer *grpc.Server
	port       int
}

func New(log *slog.Logger, port int) *App {
	gRPCServer := grpc.NewServer()
	grpcauth.Register(gRPCServer)
	return &App{
		log:        log,
		gPRCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Start() error {
	const op = "grpcapp.Start"
	log := a.log.With(slog.String("op", op))

	log.Info("Starting GRPCServer")

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("GRPCServer started", slog.String("addr", listener.Addr().String()))

	if err := a.gPRCServer.Serve(listener); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "grpcapp.Stop"
	log := a.log.With(slog.String("op", op))

	log.Info("Stopping GRPCServer")

	a.gPRCServer.GracefulStop()
}

func (a *App) MustRun() {
	if err := a.Start(); err != nil {
		panic(err)
	}
}
