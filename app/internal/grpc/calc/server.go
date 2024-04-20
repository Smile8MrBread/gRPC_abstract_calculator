package calc

import (
	"context"
	"errors"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/domain/model"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/services/calc"
	calcv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/calc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Calc interface {
	SaveExpression(ctx context.Context, expression string, ttdo, userId int64) (expressionId int64, err error)
	GetExpression(ctx context.Context, expressionId int64) (expression model.Expression, err error)
	UpdateExpression(ctx context.Context, expressionId, ttdo, result int64, status string) error
	GetAllExpressions(ctx context.Context, userId int64) (expressionsId []int64, err error)
	UpdateArithmetic(ctx context.Context, sign string, ttdo, userId int64) error
	GetArithmetic(ctx context.Context, sign string, userId int64) (arith model.Arithmetic, err error)
	AddSigns(ctx context.Context, sign string, ttdo, userId int64) (err error)
}

const EmptyValue = 0

type ServerAPI struct {
	calc Calc
	calcv1.UnimplementedCalcServer
}

func Register(grpc *grpc.Server, calc Calc) {
	calcv1.RegisterCalcServer(grpc, &ServerAPI{
		calc: calc,
	})
}

func (s *ServerAPI) SaveExpression(ctx context.Context, req *calcv1.SaveExpressionRequest) (*calcv1.SaveExpressionResponse, error) {
	if err := validateSaveExpression(req); err != nil {
		return nil, err
	}

	expressionId, err := s.calc.SaveExpression(ctx, req.GetExpression(), req.GetTtdo(), req.GetUserId())
	if err != nil {
		if errors.Is(err, calc.ErrInvalidExpression) {
			return nil, status.Error(codes.InvalidArgument, "invalid expression")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.SaveExpressionResponse{
		ExpressionId: expressionId,
	}, nil
}

func (s *ServerAPI) GetExpression(ctx context.Context, req *calcv1.GetExpressionRequest) (*calcv1.GetExpressionResponse, error) {
	if err := validateGetExpression(req); err != nil {
		return nil, err
	}

	expression, err := s.calc.GetExpression(ctx, req.GetExpressionId())
	if err != nil {
		if errors.Is(err, calc.ErrIDNotFound) {
			return nil, status.Error(codes.NotFound, "expressionId not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.GetExpressionResponse{
		Expression:   expression.Expression,
		Status:       expression.Status,
		Ttdo:         expression.Ttdo,
		ExpressionId: expression.ID,
		Result:       expression.Result,
		UserId:       expression.UserId,
	}, nil
}

func (s *ServerAPI) UpdateExpression(ctx context.Context, req *calcv1.UpdateExpressionRequest) (*calcv1.NothingMessage, error) {
	if err := validateUpdateExpression(req); err != nil {
		return nil, err
	}

	if err := s.calc.UpdateExpression(ctx, req.GetExpressionId(), req.GetTtdo(), req.GetResult(), req.GetStatus()); err != nil {
		if errors.Is(err, calc.ErrIDNotFound) {
			return nil, status.Error(codes.NotFound, "expressionId not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.NothingMessage{}, nil
}

func (s *ServerAPI) GetAllExpressions(ctx context.Context, req *calcv1.GetAllExpressionsRequest) (*calcv1.GetAllExpressionsResponse, error) {
	expressionsId, err := s.calc.GetAllExpressions(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, calc.ErrExpressionNotFound) {
			return nil, status.Error(codes.NotFound, "expression not found")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.GetAllExpressionsResponse{
		ExpressionId: expressionsId,
	}, nil
}

func (s *ServerAPI) UpdateArithmetic(ctx context.Context, req *calcv1.UpdateArithmeticRequest) (*calcv1.NothingMessage, error) {
	if err := validateUpdateArithmetic(req); err != nil {
		return nil, err
	}

	if err := s.calc.UpdateArithmetic(ctx, req.GetSign(), req.GetTtdo(), req.UserId); err != nil {
		if errors.Is(err, calc.ErrInvalidSign) {
			return nil, status.Error(codes.InvalidArgument, "invalid sign")
		}
		if errors.Is(err, calc.ErrInvalidTtdo) {
			return nil, status.Error(codes.InvalidArgument, "invalid ttdo")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.NothingMessage{}, nil
}

func (s *ServerAPI) GetArithmetic(ctx context.Context, req *calcv1.GetArithmeticRequest) (*calcv1.GetArithmeticResponse, error) {
	if err := validateGetArithmetic(req); err != nil {
		return nil, err
	}

	ttdo, err := s.calc.GetArithmetic(ctx, req.GetSign(), req.GetUserId())
	if err != nil {
		if errors.Is(err, calc.ErrInvalidSign) {
			return nil, status.Error(codes.InvalidArgument, "invalid sign")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.GetArithmeticResponse{
		Ttdo: ttdo.Ttdo,
		Sign: ttdo.Sign,
	}, nil
}

func (s *ServerAPI) AddSigns(ctx context.Context, req *calcv1.AddSignsRequest) (*calcv1.NothingMessage, error) {
	if err := validateAddSigns(req); err != nil {
		return nil, err
	}

	if err := s.calc.AddSigns(ctx, req.GetSign(), req.GetTtdo(), req.GetUserId()); err != nil {
		if errors.Is(err, calc.ErrInvalidSign) {
			return nil, status.Error(codes.InvalidArgument, "invalid sign")
		}
		if errors.Is(err, calc.ErrInvalidTtdo) {
			return nil, status.Error(codes.InvalidArgument, "invalid ttdo")
		}

		return nil, status.Error(codes.Internal, "internal error")
	}

	return &calcv1.NothingMessage{}, nil
}

func validateSaveExpression(req *calcv1.SaveExpressionRequest) error {
	if req.GetExpression() == "" {
		return status.Error(codes.InvalidArgument, "expression is empty")
	}

	if req.GetTtdo() == EmptyValue {
		return status.Error(codes.InvalidArgument, "ttdo is empty")
	}

	return nil
}

func validateGetExpression(req *calcv1.GetExpressionRequest) error {
	if req.GetExpressionId() == EmptyValue {
		return status.Error(codes.InvalidArgument, "expressionId is empty")
	}

	return nil
}

func validateUpdateExpression(req *calcv1.UpdateExpressionRequest) error {
	if req.GetExpressionId() == EmptyValue {
		return status.Error(codes.InvalidArgument, "expressionId is empty")
	}

	if req.GetStatus() == "" {
		return status.Error(codes.InvalidArgument, "status is empty")
	}

	if req.GetTtdo() == EmptyValue {
		return status.Error(codes.InvalidArgument, "ttdo is empty")
	}

	return nil
}

func validateUpdateArithmetic(req *calcv1.UpdateArithmeticRequest) error {
	if req.GetSign() == "" {
		return status.Error(codes.InvalidArgument, "sign is empty")
	}

	if req.GetTtdo() == EmptyValue {
		return status.Error(codes.InvalidArgument, "ttdo is empty")
	}

	return nil
}

func validateGetArithmetic(req *calcv1.GetArithmeticRequest) error {
	if req.GetSign() == "" {
		return status.Error(codes.InvalidArgument, "sign is empty")
	}

	return nil
}

func validateAddSigns(req *calcv1.AddSignsRequest) error {
	if req.GetSign() == "" {
		return status.Error(codes.InvalidArgument, "sign is empty")
	}

	if req.GetTtdo() == EmptyValue {
		return status.Error(codes.InvalidArgument, "ttdo is empty")
	}

	if req.GetUserId() == EmptyValue {
		return status.Error(codes.InvalidArgument, "userId is empty")
	}

	return nil
}
