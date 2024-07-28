package app

import (
	grpcapp "github.com/anizamutdinov-go/sso/internal/app/grpc"
	"log/slog"
	"time"
)

type App struct {
	GRPSCServer *grpcapp.App
}

func New(log *slog.Logger, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	//todo: init storage
	//todo: init auth layer
	//todo: init grpc-app

	grpcApp := grpcapp.New(log, grpcPort)
	return &App{GRPSCServer: grpcApp}
}
