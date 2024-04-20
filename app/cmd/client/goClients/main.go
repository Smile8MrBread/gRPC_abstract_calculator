package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/internal/services/calc"
	agentv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/agent"
	authv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/auth"
	calcv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/calc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Exp struct {
	Expression string
}

type Reg struct {
	Login    string
	Password string
}

type Arith struct {
	Sign   string
	Ttdo   int
	UserId int
}

const appSecret = "verySecretApp"

func AddSigns(calcClient calcv1.CalcClient, ttdo int64, claims jwt.MapClaims, log *slog.Logger) error {
	_, err := calcClient.AddSigns(context.Background(), &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "+",
		Ttdo:   ttdo,
	})
	_, err = calcClient.AddSigns(context.Background(), &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "*",
		Ttdo:   ttdo,
	})
	_, err = calcClient.AddSigns(context.Background(), &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "-",
		Ttdo:   ttdo,
	})
	_, err = calcClient.AddSigns(context.Background(), &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "/",
		Ttdo:   ttdo,
	})

	if err != nil {
		log.Error("Failed to create signs")
		return err
	}

	return nil
}

func CheckToken(r *http.Request) bool {
	cookie, err := r.Cookie("token")
	if err != nil || cookie.Value == "" {
		return false
	}

	return true
}

func ParseCookie(r *http.Request, log *slog.Logger) (jwt.MapClaims, error) {
	token, err := r.Cookie("token")
	if err != nil {
		log.Error("Error getting cookie", err.Error())
		return nil, err
	}

	tokenParsed, err := jwt.Parse(token.Value, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	if !ok {
		log.Error("Error parsing token", err.Error())
		return nil, err
	}

	return claims, nil
}

func GetAllArithmetics(calcClient calcv1.CalcClient, userId int64) map[string]int64 {
	respPlus, err := calcClient.GetArithmetic(context.Background(), &calcv1.GetArithmeticRequest{
		Sign:   "+",
		UserId: userId,
	})
	if err != nil {
		fmt.Errorf("%v", err)
		return nil
	}

	respMinus, err := calcClient.GetArithmetic(context.Background(), &calcv1.GetArithmeticRequest{
		Sign:   "-",
		UserId: userId,
	})
	if err != nil {
		fmt.Errorf("%v", err)
		return nil
	}

	respDev, err := calcClient.GetArithmetic(context.Background(), &calcv1.GetArithmeticRequest{
		Sign:   "/",
		UserId: userId,
	})
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return nil
	}

	respMult, err := calcClient.GetArithmetic(context.Background(), &calcv1.GetArithmeticRequest{
		Sign:   "*",
		UserId: userId,
	})
	if err != nil {
		_ = fmt.Errorf("%v", err)
		return nil
	}

	result := map[string]int64{
		respPlus.Sign:  respPlus.Ttdo,
		respMinus.Sign: respMinus.Ttdo,
		respDev.Sign:   respDev.Ttdo,
		respMult.Sign:  respMult.Ttdo,
	}

	return result
}
func main() {
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	conn, err := grpc.NewClient(":26026", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	authClient := authv1.NewAuthClient(conn)

	conn1, err := grpc.NewClient(":26126", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	calcClient := calcv1.NewCalcClient(conn1)

	conn2, err := grpc.NewClient(":26226", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	if err != nil {
		panic(err)
	}
	agentClient := agentv1.NewGRPCAgentClient(conn2)

	// Reg and Auth page
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("app/cmd/client/html/index.html")
		if err != nil {
			log.Error("Error parsing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Error("Error executing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})

	// Auth service
	// Handler for recording a new user
	r.Post("/register", func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)

		data := Reg{}
		err := body.Decode(&data)
		if err != nil {
			log.Error("Error parsing json")
			return
		}

		_, err = authClient.Register(r.Context(), &authv1.RegisterRequest{
			Login:    data.Login,
			Password: data.Password,
		})
		if err != nil {
			if errors.Is(err, status.Error(codes.AlreadyExists, "user already exists")) {
				w.WriteHeader(http.StatusConflict)
				log.Error("User already exists")
			}

			log.Error("Error registering ", err.Error())
			return
		}
	})
	// Handler for login and extradition new jwt token's and user's id in cookies
	r.Post("/login", func(w http.ResponseWriter, r *http.Request) {
		body := json.NewDecoder(r.Body)

		data := Reg{}
		err := body.Decode(&data)
		if err != nil {
			log.Error("Error parsing json")
			return
		}

		token, err := authClient.Login(r.Context(), &authv1.LoginRequest{
			Login:    data.Login,
			Password: data.Password,
			AppId:    1,
		})
		if err != nil {
			if errors.Is(err, status.Error(codes.InvalidArgument, "invalid login or password")) {
				w.WriteHeader(http.StatusNotAcceptable)
				log.Error("Invalid login or password")
			}

			log.Error("Error logging in ", err.Error())
			return
		}

		tokenParsed, err := jwt.Parse(token.Token, func(token *jwt.Token) (interface{}, error) {
			return []byte(appSecret), nil
		})

		claims, ok := tokenParsed.Claims.(jwt.MapClaims)
		if !ok {
			log.Error("Error parsing token", err.Error())
			return
		}

		_, err = calcClient.GetArithmetic(context.Background(), &calcv1.GetArithmeticRequest{
			Sign:   "+",
			UserId: int64(claims["user_id"].(float64)),
		})
		if err != nil {
			err = AddSigns(calcClient, 1, claims, log)
			if err != nil {
				log.Error("Error adding signs", err.Error())
				return
			}
		}

		id := strconv.Itoa(int(claims["user_id"].(float64)))

		http.SetCookie(w, &http.Cookie{
			Name:   "token",
			Value:  token.Token,
			Path:   "/",
			Secure: false,
		})

		http.SetCookie(w, &http.Cookie{
			Name:   "id",
			Value:  id,
			Path:   "/",
			Secure: false,
		})

		http.Redirect(w, r, "/newExp", http.StatusSeeOther)
	})

	// Calculator service (orkestrator and agent)
	// Page with form to add a new expression
	r.Get("/newExp", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Warn("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		ts, err := template.ParseFiles("app/cmd/client/html/newExp.html")
		if err != nil {
			log.Error("Error parsing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Error("Error executing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	// Handler for adding a new expression
	r.Post("/newExp", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Error("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
		}

		body := json.NewDecoder(r.Body)

		data := Exp{}

		err := body.Decode(&data)
		if err != nil {
			log.Error("Error decoding body ", err.Error())
			return
		}

		limitExpression := regexp.MustCompile(`^-?[0-9]+(?:[+\-*/][0-9]+)*$`)
		if !limitExpression.MatchString(data.Expression) {
			log.Warn("Invalid expression", slog.String("error", calc.ErrInvalidExpression.Error()), data.Expression)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		claims, err := ParseCookie(r, log)
		if err != nil {
			log.Error("Error parsing cookie", err.Error())
			return
		}

		arithTime := GetAllArithmetics(calcClient, int64(claims["user_id"].(float64)))

		ttdo := int64(strings.Count(data.Expression, "-"))*arithTime["-"] + int64(strings.Count(data.Expression, "+"))*arithTime["+"]
		ttdo += int64(strings.Count(data.Expression, "/"))*arithTime["/"] + int64(strings.Count(data.Expression, "*"))*arithTime["*"]

		exp := claims["exp"].(float64)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(int64(exp)-time.Now().Unix())*10e8)

		resp, err := calcClient.SaveExpression(ctx, &calcv1.SaveExpressionRequest{
			Expression: data.Expression,
			Ttdo:       ttdo,
			UserId:     int64(claims["user_id"].(float64)),
		})
		if err != nil {
			if errors.Is(err, status.Error(codes.DeadlineExceeded, "context deadline exceeded")) {
				http.SetCookie(w, &http.Cookie{
					Name:   "token",
					Value:  "",
					Path:   "/",
					Secure: false,
				})
				http.SetCookie(w, &http.Cookie{
					Name:   "id",
					Value:  "",
					Path:   "/",
					Secure: false,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			if errors.Is(err, status.Error(codes.InvalidArgument, "ttdo is empty")) {
				http.Error(w, "Ttdo is empty", http.StatusBadRequest)
				log.Error("Ttdo is empty")
			}

			log.Error("Error saving expression ", err.Error())
			return
		}

		go func() {
			_, err = agentClient.ExpForDo(ctx, &agentv1.ExpForDoRequest{
				ExpressionId: resp.ExpressionId,
			})
			if err != nil {
				log.Error("Error saving expression ", err.Error())
				return
			}
		}()

	})
	// Results page
	r.Get("/allExp", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Warn("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		ts, err := template.ParseFiles("app/cmd/client/html/allExp.html")
		if err != nil {
			log.Error("Error parsing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Error("Error executing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	// Handler for getting all expressions by user
	r.Get("/allExp/{userId}", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Warn("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		claims, err := ParseCookie(r, log)
		if err != nil {
			log.Error("Error parsing cookie", err.Error())
			return
		}

		exp := claims["exp"].(float64)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(int64(exp)-time.Now().Unix())*10e8)

		userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil {
			log.Error("Error parsing userId", err.Error())
			return
		}
		ids, err := calcClient.GetAllExpressions(ctx, &calcv1.GetAllExpressionsRequest{
			UserId: int64(userId),
		})
		if err != nil {
			if errors.Is(err, status.Error(codes.DeadlineExceeded, "context deadline exceeded")) {
				http.SetCookie(w, &http.Cookie{
					Name:   "token",
					Value:  "",
					Path:   "/",
					Secure: false,
				})
				http.SetCookie(w, &http.Cookie{
					Name:   "id",
					Value:  "",
					Path:   "/",
					Secure: false,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			log.Error("Error getting all expressions ", err.Error())
			return
		}

		var res []map[string]string
		for i := range ids.ExpressionId {
			exp, err := calcClient.GetExpression(ctx, &calcv1.GetExpressionRequest{
				ExpressionId: ids.ExpressionId[i],
			})
			if err != nil {
				log.Error("Error getting expression", err.Error())
				return
			}
			res = append(res, map[string]string{
				"expressionId": strconv.FormatInt(exp.ExpressionId, 10),
				"expression":   exp.Expression,
				"status":       exp.Status,
				"result":       strconv.FormatInt(exp.Result, 10),
			})
		}

		data, err := json.Marshal(res)
		if err != nil {
			log.Error("Error marshaling json", err.Error())
			return
		}

		_, err = w.Write(data)
		if err != nil {
			log.Error("Error writing response", err.Error())
			return
		}

		w.Header().Set("Content-Type", "application/json")
		return
	})
	// Page with forms to set time to do to each sign
	r.Get("/ariths", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Warn("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		ts, err := template.ParseFiles("app/cmd/client/html/ariths.html")
		if err != nil {
			log.Error("Error parsing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err = ts.Execute(w, nil)
		if err != nil {
			log.Error("Error executing template")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	})
	// Handler for getting all time to do for each sign by user
	r.Get("/ariths/{userId}", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Warn("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		userId, err := strconv.Atoi(chi.URLParam(r, "userId"))
		if err != nil {
			log.Error("Error parsing userId", err.Error())
			return
		}
		allSigns := GetAllArithmetics(calcClient, int64(userId))

		data, err := json.Marshal(allSigns)
		if err != nil {
			log.Error("Error marshaling json", err.Error())
			return
		}

		w.WriteHeader(http.StatusCreated)
		w.Write(data)

		w.Header().Set("Content-Type", "application/json")
		return
	})
	// Handle to update sign's time of this user
	r.Patch("/ariths/{userId}", func(w http.ResponseWriter, r *http.Request) {
		if !CheckToken(r) {
			log.Warn("Token not found")

			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		claims, err := ParseCookie(r, log)
		if err != nil {
			log.Error("Error parsing cookie", err.Error())
			return
		}

		body := json.NewDecoder(r.Body)

		data := Arith{}

		err = body.Decode(&data)
		if err != nil {
			log.Error("Error decoding body ", err.Error())
			return
		}

		exp := claims["exp"].(float64)
		ctx, _ := context.WithTimeout(context.Background(), time.Duration(int64(exp)-time.Now().Unix())*10e8)

		limitTTDO := regexp.MustCompile(`^[1-9][0-9]*?$`)
		if !limitTTDO.MatchString(strconv.Itoa(data.Ttdo)) {
			w.WriteHeader(http.StatusBadRequest)
			log.Error("invalid ttdo response")
			return
		}

		_, err = calcClient.UpdateArithmetic(ctx, &calcv1.UpdateArithmeticRequest{
			UserId: int64(data.UserId),
			Ttdo:   int64(data.Ttdo),
			Sign:   data.Sign,
		})
		if err != nil {
			if errors.Is(err, status.Error(codes.DeadlineExceeded, "context deadline exceeded")) {
				http.SetCookie(w, &http.Cookie{
					Name:   "token",
					Value:  "",
					Path:   "/",
					Secure: false,
				})
				http.SetCookie(w, &http.Cookie{
					Name:   "id",
					Value:  "",
					Path:   "/",
					Secure: false,
				})
				http.Redirect(w, r, "/", http.StatusSeeOther)
				return
			}

			log.Error("update arithmetic request error", err)
			return
		}
	})

	fmt.Println("Server is running on port 8080")
	err = http.ListenAndServe(":8080", r)
	if err != nil {
		panic(err.Error())
	}
}
