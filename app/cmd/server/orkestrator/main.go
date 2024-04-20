package main

import (
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/app"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/config"
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
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting orkestrator gRPC server...")

	application := app.New(log, cfg.AuthGRPC.Port, cfg.CalcGRPC.Port, cfg.AgentGRPC.Port, cfg.StoragePath, cfg.TokenTTL)

	go application.GRPCSrv.MustRunCalc()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	log.Info("Stopping orkestrator gRPC server...",
		slog.Any("signal", stopSign),
	)

	application.GRPCSrv.StopCalc()
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
