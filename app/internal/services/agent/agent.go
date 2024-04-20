package agent

import (
	"context"
	"fmt"
	"log/slog"
)

type AgentSender interface {
	ExpForDo(ctx context.Context, expressionId int64) error
}

type Agent struct {
	log    *slog.Logger
	sender AgentSender
}

func New(log *slog.Logger, sender AgentSender) *Agent {
	return &Agent{
		log:    log,
		sender: sender,
	}
}

func (a *Agent) ExpForDo(ctx context.Context, expressionId int64) error {
	const op = "agent.SendExpression"

	log := a.log.With(slog.String("op", op), slog.Int64("expressionId", expressionId))

	err := a.sender.ExpForDo(ctx, expressionId)
	if err != nil {
		log.Warn("Error sending expression", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
