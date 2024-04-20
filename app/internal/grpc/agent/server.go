package agent

import (
	"context"
	agentv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/agent"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Agent interface {
	ExpForDo(ctx context.Context, expressionId int64) error
}

type ServerAPI struct {
	agentv1.UnimplementedGRPCAgentServer
	agent Agent
}

const emptyValue = 0

func Register(gRPC *grpc.Server, agent Agent) {
	agentv1.RegisterGRPCAgentServer(gRPC, &ServerAPI{agent: agent})
}

func (s *ServerAPI) ExpForDo(ctx context.Context, req *agentv1.ExpForDoRequest) (*agentv1.NothingMessage, error) {
	if err := validateExpForDo(req); err != nil {
		return nil, err
	}

	if err := s.agent.ExpForDo(ctx, req.GetExpressionId()); err != nil {
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &agentv1.NothingMessage{}, nil
}

func validateExpForDo(req *agentv1.ExpForDoRequest) error {
	if req.GetExpressionId() == emptyValue {
		return status.Error(codes.InvalidArgument, "expressionId is empty")
	}

	return nil
}
