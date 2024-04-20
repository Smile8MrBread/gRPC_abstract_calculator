package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/domain/model"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/storage"
	"github.com/mattn/go-sqlite3"
	"regexp"
	"slices"
	"strconv"
	"sync"
	"time"
)

type Storage struct {
	db *sql.DB
}

func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) Stop() error {
	return s.db.Close()
}

// SaveUser saves user to db.
func (s *Storage) SaveUser(ctx context.Context, login string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	stmt, err := s.db.Prepare("INSERT INTO users(login, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, login, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by login.
func (s *Storage) User(ctx context.Context, login string) (model.User, error) {
	const op = "storage.sqlite.User"

	stmt, err := s.db.Prepare("SELECT id, login, pass_hash FROM users WHERE login = ?")
	if err != nil {
		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, login)

	var user model.User
	err = row.Scan(&user.ID, &user.Login, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}

		return model.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

// App returns app by id.
func (s *Storage) App(ctx context.Context, id int) (model.App, error) {
	const op = "storage.sqlite.App"

	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, id)

	var app model.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}

		return model.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}

func (s *Storage) SaveExpression(ctx context.Context, expression string, ttdo, userId int64) (int64, error) {
	const op = "storage.sqlite.SaveExpression"

	stmt, err := s.db.Prepare("INSERT INTO expressions(expression, ttdo, status, user_id) VALUES(?,?,?,?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, expression, ttdo, "InProcess", userId)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrExpressionExists)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	expressionId, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return expressionId, nil
}

func (s *Storage) GetAllExpressions(ctx context.Context, userId int64) ([]int64, error) {
	const op = "storage.sqlite.GetExpressions"

	rows, err := s.db.QueryContext(ctx, "SELECT id FROM expressions WHERE user_id = ?", userId)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var expressionsId []int64
	for rows.Next() {
		var expressionId int64
		if err := rows.Scan(&expressionId); err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, fmt.Errorf("%s: %w", op, storage.ErrExceptionNotFound)
			}

			return nil, fmt.Errorf("%s: %w", op, err)
		}
		expressionsId = append(expressionsId, expressionId)
	}

	return expressionsId, nil
}

func (s *Storage) GetExpression(ctx context.Context, expressionId int64) (model.Expression, error) {
	const op = "storage.sqlite.GetExpression"

	stmt, err := s.db.Prepare("SELECT id, expression, ttdo, status, result, user_id FROM expressions WHERE id =?")
	if err != nil {
		return model.Expression{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, expressionId)

	var expression model.Expression
	err = row.Scan(&expression.ID, &expression.Expression, &expression.Ttdo, &expression.Status, &expression.Result, &expression.UserId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Expression{}, fmt.Errorf("%s: %w", op, storage.ErrExceptionNotFound)
		}

		return model.Expression{}, fmt.Errorf("%s: %w", op, err)
	}

	return expression, nil
}

func (s *Storage) UpdateExpression(ctx context.Context, expressionId, ttdo, result int64, status string) error {
	const op = "storage.sqlite.UpdateExpression"

	stmt, err := s.db.Prepare("UPDATE expressions SET ttdo =?, status =?, result=? WHERE id =?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, ttdo, status, result, expressionId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// GetArithmetic return the model.arithmetic of given sign
func (s *Storage) GetArithmetic(ctx context.Context, sign string, userId int64) (model.Arithmetic, error) {
	const op = "storage.sqlite.GetArithmetic"

	stmt, err := s.db.Prepare("SELECT sign, ttdo FROM arithmetics WHERE sign =? AND user_id =?")
	if err != nil {
		return model.Arithmetic{}, fmt.Errorf("%s: %w", op, err)
	}

	row := stmt.QueryRowContext(ctx, sign, userId)

	var arithmetic model.Arithmetic
	err = row.Scan(&arithmetic.Sign, &arithmetic.Ttdo)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Arithmetic{}, fmt.Errorf("%s: %w", op, storage.ErrInvalidSign)
		}

		return model.Arithmetic{}, fmt.Errorf("%s: %w", op, err)
	}

	return arithmetic, nil
}

func (s *Storage) UpdateArithmetic(ctx context.Context, sign string, ttdo, userId int64) error {
	const op = "storage.sqlite.UpdateArithmetic"

	stmt, err := s.db.Prepare("UPDATE arithmetics SET ttdo =? WHERE sign =? AND user_id =?")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, ttdo, sign, userId)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// AddSigns add all signs for new user
func (s *Storage) AddSigns(ctx context.Context, sign string, ttdo, userId int64) error {
	const op = "storage.sqlite.AddSigns"

	stmt, err := s.db.Prepare("INSERT INTO arithmetics(sign, ttdo, user_id) VALUES(?,?,?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.ExecContext(ctx, sign, ttdo, userId)
	if err != nil {
		var sqliteErr sqlite3.Error
		if errors.As(err, &sqliteErr) && errors.Is(sqliteErr.ExtendedCode, sqlite3.ErrConstraintUnique) {
			return fmt.Errorf("%s: %w", op, storage.ErrSignExists)
		}

		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// ExpForDo keep exp id and in go routine image hard work about given exp
func (s *Storage) ExpForDo(ctx context.Context, expressionId int64) error {
	const op = "storage.sqlite.ExpForDo"
	var example []string

	wg := &sync.WaitGroup{}
	wg.Add(1)
	var errGl error
	go func() {
		defer wg.Done()

		exp, err := s.GetExpression(ctx, expressionId)
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}
		timeDiv, err := s.GetArithmetic(ctx, "/", exp.UserId)
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}

		timeMult, err := s.GetArithmetic(ctx, "*", exp.UserId)
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}

		timeMinus, err := s.GetArithmetic(ctx, "-", exp.UserId)
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}

		timePlus, err := s.GetArithmetic(ctx, "+", exp.UserId)
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}

		re := regexp.MustCompile(`\d+|\+|-|\*|/`)
		subExps := re.FindAllString(exp.Expression, -1)
		example = subExps
		fmt.Println("I need to do:", subExps)

		for {
			indDiv := slices.Index(subExps, "/")
			indMult := slices.Index(subExps, "*")

			if indDiv == -1 && indMult == -1 {
				break
			}

			if indDiv != -1 && (indDiv < indMult || indMult == -1) {

				fmt.Println(expressionId, "=> I do", subExps[indDiv-1], "/", subExps[indDiv+1], "in expression:", subExps, "of user:", exp.UserId)
				time.Sleep(time.Duration(timeDiv.Ttdo) * time.Second)

				a, _ := strconv.ParseFloat(subExps[indDiv-1], 64)
				b, _ := strconv.ParseFloat(subExps[indDiv+1], 64)
				subExps = slices.Delete(subExps, indDiv-1, indDiv+2)

				arr := subExps[:indDiv-1]
				brr := subExps[indDiv-1:]

				c := make([]string, len(arr)+len(brr)+1)
				copy(c[:len(arr)], arr)
				copy(c[len(arr):], []string{fmt.Sprintf("%.2f", a/b)})
				copy(c[len(arr)+1:], brr)
				subExps = c

			} else {

				fmt.Println(expressionId, "=> I do", subExps[indMult-1], "*", subExps[indMult+1], "in expression:", subExps, "of user:", exp.UserId)
				time.Sleep(time.Duration(timeMult.Ttdo) * time.Second)

				a, _ := strconv.ParseFloat(subExps[indMult-1], 64)
				b, _ := strconv.ParseFloat(subExps[indMult+1], 64)
				subExps = slices.Delete(subExps, indMult-1, indMult+2)

				arr := subExps[:indMult-1]
				brr := subExps[indMult-1:]

				c := make([]string, len(arr)+len(brr)+1)
				copy(c[:len(arr)], arr)
				copy(c[len(arr):], []string{fmt.Sprintf("%.2f", a*b)})
				copy(c[len(arr)+1:], brr)
				subExps = c
			}

		}
		for {
			indPlus := slices.Index(subExps, "+")
			indMinus := slices.Index(subExps, "-")

			if indPlus == -1 && indMinus == -1 {
				break
			}
			if indMinus != -1 && (indMinus < indPlus || indPlus == -1) {
				fmt.Println(expressionId, "=> I do", subExps[indMinus-1], "-", subExps[indMinus+1], "in expression:", subExps, "of user:", exp.UserId)
				time.Sleep(time.Duration(timeMinus.Ttdo) * time.Second)

				a, _ := strconv.ParseFloat(subExps[indMinus-1], 64)
				b, _ := strconv.ParseFloat(subExps[indMinus+1], 64)
				subExps = slices.Delete(subExps, indMinus-1, indMinus+2)

				arr := subExps[:indMinus-1]
				brr := subExps[indMinus-1:]

				c := make([]string, len(arr)+len(brr)+1)
				copy(c[:len(arr)], arr)
				copy(c[len(arr):], []string{fmt.Sprintf("%.2f", a-b)})
				copy(c[len(arr)+1:], brr)
				subExps = c

			} else {
				fmt.Println(expressionId, "=> I do", subExps[indPlus-1], "+", subExps[indPlus+1], "in expression:", subExps, "of user:", exp.UserId)
				time.Sleep(time.Duration(timePlus.Ttdo) * time.Second)

				a, _ := strconv.ParseFloat(subExps[indPlus-1], 64)
				b, _ := strconv.ParseFloat(subExps[indPlus+1], 64)
				subExps = slices.Delete(subExps, indPlus-1, indPlus+2)

				arr := subExps[:indPlus-1]
				brr := subExps[indPlus-1:]

				c := make([]string, len(arr)+len(brr)+1)
				copy(c[:len(arr)], arr)
				copy(c[len(arr):], []string{fmt.Sprintf("%.2f", a+b)})
				copy(c[len(arr)+1:], brr)
				subExps = c
			}
		}

		res, err := strconv.ParseFloat(subExps[0], 64)
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}

		err = s.UpdateExpression(ctx, exp.ID, 0, int64(res), "Done")
		if err != nil {
			errGl = fmt.Errorf("%s: %w", op, err)
			return
		}

		fmt.Println(expressionId, "=> I done the expression:", example, " and result is:", subExps[0])
	}()

	wg.Wait()

	if errGl != nil {
		return errGl
	}

	return nil
}
