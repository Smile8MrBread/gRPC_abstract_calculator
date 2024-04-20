package suite

import (
	"context"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/config"
	agentv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/agent"
	ssov1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/auth"
	calcv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
	"strconv"
	"testing"
)

const grpcHost = "localhost"

type Suite struct {
	*testing.T
	Cfg         *config.Config
	AuthClient  ssov1.AuthClient
	CalcClient  calcv1.CalcClient
	AgentClient agentv1.GRPCAgentClient
}

func NewSuite(t *testing.T) (context.Context, *Suite) {
	t.Helper()
	t.Parallel()

	cfg := config.MustLoadPath("../config/local.yaml")

	ctx, cansel := context.WithTimeout(context.Background(), cfg.AuthGRPC.Timeout)

	t.Cleanup(func() {
		t.Helper()
		cansel()
	})

	cc, err := grpc.NewClient(net.JoinHostPort(grpcHost, strconv.Itoa(cfg.AuthGRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial gRPC server: %s", err)
	}

	cc1, err := grpc.NewClient(net.JoinHostPort(grpcHost, strconv.Itoa(cfg.CalcGRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial gRPC server: %s", err)
	}

	cc2, err := grpc.NewClient(net.JoinHostPort(grpcHost, strconv.Itoa(cfg.AgentGRPC.Port)),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("failed to dial gRPC server: %s", err)
	}

	return ctx, &Suite{
		T:           t,
		Cfg:         cfg,
		AuthClient:  ssov1.NewAuthClient(cc),
		CalcClient:  calcv1.NewCalcClient(cc1),
		AgentClient: agentv1.NewGRPCAgentClient(cc2),
	}
}
