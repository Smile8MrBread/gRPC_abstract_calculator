package calc

import (
	"context"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/domain/model"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/storage"
	"log/slog"
)

var (
	ErrExpressionNotFound = errors.New("expression not found")
	ErrInvalidExpression  = errors.New("invalid expression")
	ErrInvalidTtdo        = errors.New("invalid ttdo")
	ErrIDNotFound         = errors.New("id not found")
	ErrInvalidSign        = errors.New("invalid sign")
)

type Calc struct {
	log      *slog.Logger
	saver    Saver
	provider Provider
	updater  Updater
}

type Saver interface {
	SaveExpression(ctx context.Context, expression string, ttdo, userId int64) (int64, error)
	GetAllExpressions(ctx context.Context, userId int64) (expressionsId []int64, err error)
	AddSigns(ctx context.Context, sign string, ttdo, userId int64) (err error)
}

type Provider interface {
	GetExpression(ctx context.Context, expressionId int64) (expression model.Expression, err error)
	GetArithmetic(ctx context.Context, sign string, userId int64) (arith model.Arithmetic, err error)
}

type Updater interface {
	UpdateExpression(ctx context.Context, expressionId, ttdo, result int64, status string) error
	UpdateArithmetic(ctx context.Context, sign string, ttdo, userId int64) error
}

func NewCalc(log *slog.Logger, saver Saver, provider Provider, updater Updater) *Calc {
	return &Calc{
		log:      log,
		saver:    saver,
		provider: provider,
		updater:  updater,
	}
}

func (c *Calc) GetExpression(ctx context.Context, expressionId int64) (expression model.Expression, err error) {
	const op = "calc.GetExpression"

	log := c.log.With(slog.String("op", op), slog.Int64("expressionId", expressionId))
	log.Info("Getting expression")

	expression, err = c.provider.GetExpression(ctx, expressionId)
	if err != nil {
		if errors.Is(err, ErrExpressionNotFound) {
			log.Warn("Expression not found", slog.String("error", err.Error()))
			return model.Expression{}, ErrExpressionNotFound
		}

		log.Warn("Error getting expression", slog.String("error", err.Error()))
		return model.Expression{}, fmt.Errorf("%s: %w", op, err)
	}

	return expression, nil
}

func (c *Calc) GetArithmetic(ctx context.Context, sign string, userId int64) (arith model.Arithmetic, err error) {
	const op = "calc.GetArithmetic"

	log := c.log.With(slog.String("op", op), slog.String("sign", sign))
	log.Info("Getting arithmetic")

	arith, err = c.provider.GetArithmetic(ctx, sign, userId)
	if err != nil {
		if errors.Is(err, ErrInvalidSign) {
			log.Warn("Invalid sign", slog.String("error", err.Error()))
			return model.Arithmetic{}, ErrInvalidSign
		}

		log.Warn("Error getting arithmetic", slog.String("error", err.Error()))
		return model.Arithmetic{}, fmt.Errorf("%s: %w", op, err)
	}

	return arith, nil
}

func (c *Calc) UpdateExpression(ctx context.Context, expressionId, ttdo, result int64, status string) error {
	const op = "calc.UpdateExpression"

	log := c.log.With(slog.String("op", op), slog.Int64("expressionId", expressionId))
	log.Info("Updating expression")

	err := c.updater.UpdateExpression(ctx, expressionId, ttdo, result, status)
	if err != nil {
		log.Warn("Error updating expression", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Calc) UpdateArithmetic(ctx context.Context, sign string, ttdo, userId int64) error {
	const op = "calc.UpdateArithmetic"

	log := c.log.With(slog.String("op", op), slog.String("sign", sign))
	log.Info("Updating arithmetic")

	err := c.updater.UpdateArithmetic(ctx, sign, ttdo, userId)
	if err != nil {
		log.Warn("Error updating arithmetic", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (c *Calc) SaveExpression(ctx context.Context, expression string, ttdo, userId int64) (int64, error) {
	const op = "calc.SaveExpression"

	log := c.log.With(slog.String("op", op))
	log.Info("Saving expression")

	expressionId, err := c.saver.SaveExpression(ctx, expression, ttdo, userId)
	if err != nil {
		if errors.Is(err, storage.ErrExpressionExists) {
			log.Warn("Expression exists", slog.String("error", err.Error()))
			return 0, storage.ErrExpressionExists
		}

		log.Warn("Error saving expression", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return expressionId, nil
}

func (c *Calc) GetAllExpressions(ctx context.Context, userId int64) (expressionsId []int64, err error) {
	const op = "calc.GetAllExpressions"

	log := c.log.With(slog.String("op", op))
	log.Info("Getting all expressions")

	expressionsId, err = c.saver.GetAllExpressions(ctx, userId)
	if err != nil {
		if errors.Is(err, storage.ErrExpressionExists) {
			log.Warn("Expression exists", slog.String("error", err.Error()))
			return nil, storage.ErrExceptionNotFound
		}

		log.Warn("Error getting all expressions", slog.String("error", err.Error()))
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return expressionsId, nil
}

func (c *Calc) AddSigns(ctx context.Context, sign string, ttdo, userId int64) (err error) {
	const op = "calc.AddSigns"

	log := c.log.With(slog.String("op", op), slog.String("sign", sign))
	log.Info("Adding signs")

	err = c.saver.AddSigns(ctx, sign, ttdo, userId)
	if err != nil {
		log.Warn("Error adding signs", slog.String("error", err.Error()))
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}
