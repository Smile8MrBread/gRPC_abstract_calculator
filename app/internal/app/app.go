package app

import (
	grpcapp "github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/app/grpc"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/services/agent"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/services/auth"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/services/calc"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/storage/sqlite"
	"log/slog"
	"time"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	authPort int,
	calcPort int,
	agentPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	agentService := agent.New(log, storage)

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	calcService := calc.NewCalc(log, storage, storage, storage)

	grpcApp := grpcapp.New(log, authService, calcService, agentService, agentPort, authPort, calcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
