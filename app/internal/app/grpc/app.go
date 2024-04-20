package grpcapp

import (
	"fmt"
	agentgrpc "github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/grpc/agent"
	authgrpc "github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/grpc/auth"
	calcgrpc "github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/grpc/calc"
	"google.golang.org/grpc"
	"log/slog"
	"net"
	"strconv"
)

type App struct {
	log       *slog.Logger
	auth      *grpc.Server
	calc      *grpc.Server
	agent     *grpc.Server
	agentPort int
	authPort  int
	calcPort  int
}

func New(
	log *slog.Logger,
	authServ authgrpc.Auth,
	calcServ calcgrpc.Calc,
	agentServ agentgrpc.Agent,
	agentPort int,
	authPort int,
	calcPort int,
) *App {
	gRPCServer := grpc.NewServer()

	agentgrpc.Register(gRPCServer, agentServ)

	calcgrpc.Register(gRPCServer, calcServ)

	authgrpc.Register(gRPCServer, authServ)

	return &App{
		log:       log,
		auth:      gRPCServer,
		calc:      gRPCServer,
		agent:     gRPCServer,
		agentPort: agentPort,
		authPort:  authPort,
		calcPort:  calcPort,
	}
}

func (a *App) MustRunAuth() {
	if err := a.RunAuth(); err != nil {
		panic(err)
	}
}

func (a *App) RunAuth() error {
	const op = "gRPC.RunAuth"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.authPort),
	)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.authPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("auth server is running", slog.String("address", l.Addr().String()))

	if err := a.auth.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) StopAuth() {
	const op = "gRPC.StopAuth"

	a.log.With(slog.String("op", op)).
		Info("auth server stopping", slog.Int("port", a.authPort))

	a.auth.GracefulStop()
}

func (a *App) MustRunCalc() {
	if err := a.RunCalc(); err != nil {
		panic(err)
	}
}

func (a *App) RunCalc() error {
	const op = "gRPC.RunCalc"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.calcPort),
	)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.calcPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("calc server is running", slog.String("address", l.Addr().String()))

	if err := a.calc.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) StopCalc() {
	const op = "gRPC.StopCalc"

	a.log.With(slog.String("op", op)).
		Info("calc server stopping", slog.Int("port", a.calcPort))

	a.calc.GracefulStop()
}

func (a *App) MustRunAgent() {
	if err := a.RunAgent(); err != nil {
		panic(err)
	}
}

func (a *App) RunAgent() error {
	const op = "gRPC.RunAgent"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.agentPort),
	)

	l, err := net.Listen("tcp", ":"+strconv.Itoa(a.agentPort))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("agent server is running", slog.String("address", l.Addr().String()))

	if err := a.agent.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) StopAgent() {
	const op = "gRPC.StopAgent"

	a.log.With(slog.String("op", op)).
		Info("agent server stopping", slog.Int("port", a.agentPort))

	a.agent.GracefulStop()
}
