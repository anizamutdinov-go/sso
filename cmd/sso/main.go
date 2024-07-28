package main

import (
	"github.com/anizamutdinov-go/sso/internal/app"
	"github.com/anizamutdinov-go/sso/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {

	cfg := config.Mustload()
	log := setupLogger(cfg.Env)

	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPSCServer.MustRun()

	//todo: engine
	//todo: start grpc-app
	//todo: ???

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	stopSign := <-stop
	log.Info("Stopping SSO-application", slog.String("signal", stopSign.String()))
	application.GRPSCServer.Stop()
}

func setupLogger(env string) (log *slog.Logger) {
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
	}
	return
}
