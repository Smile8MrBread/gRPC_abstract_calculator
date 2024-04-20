package tests

import (
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/tests/suite"
	ssov1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/auth"
	calcv1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/calc"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCalc_Happy(t *testing.T) {
	ctx, su := suite.NewSuite(t)

	email := gofakeit.Email()

	password := randomPassword()

	_, err := su.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    email,
		Password: password,
	})

	respLogin, err := su.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    email,
		Password: password,
		AppId:    appID,
	})

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	_, err = su.CalcClient.AddSigns(ctx, &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "-",
		Ttdo:   1,
	})
	require.NoError(t, err)

	_, err = su.CalcClient.AddSigns(ctx, &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "+",
		Ttdo:   5,
	})
	require.NoError(t, err)

	_, err = su.CalcClient.AddSigns(ctx, &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "/",
		Ttdo:   2,
	})
	require.NoError(t, err)

	_, err = su.CalcClient.AddSigns(ctx, &calcv1.AddSignsRequest{
		UserId: int64(claims["user_id"].(float64)),
		Sign:   "*",
		Ttdo:   3,
	})
	require.NoError(t, err)

	plus, err := su.CalcClient.GetArithmetic(ctx, &calcv1.GetArithmeticRequest{
		Sign:   "+",
		UserId: int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	minus, err := su.CalcClient.GetArithmetic(ctx, &calcv1.GetArithmeticRequest{
		Sign:   "-",
		UserId: int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	mult, err := su.CalcClient.GetArithmetic(ctx, &calcv1.GetArithmeticRequest{
		Sign:   "*",
		UserId: int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	div, err := su.CalcClient.GetArithmetic(ctx, &calcv1.GetArithmeticRequest{
		Sign:   "/",
		UserId: int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	require.Equal(t, minus.Ttdo, int64(1))
	require.Equal(t, plus.Ttdo, int64(5))
	require.Equal(t, div.Ttdo, int64(2))
	require.Equal(t, mult.Ttdo, int64(3))

	expId, err := su.CalcClient.SaveExpression(ctx, &calcv1.SaveExpressionRequest{
		Expression: "2-2+3/4*3",
		Ttdo:       999,
		UserId:     int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	exp, err := su.CalcClient.GetExpression(ctx, &calcv1.GetExpressionRequest{
		ExpressionId: expId.ExpressionId,
	})
	require.NoError(t, err)

	require.Equal(t, exp.Expression, "2-2+3/4*3")
	require.Equal(t, exp.Ttdo, int64(999))
	require.Equal(t, int64(claims["user_id"].(float64)), exp.UserId)

	_, err = su.CalcClient.UpdateExpression(ctx, &calcv1.UpdateExpressionRequest{
		ExpressionId: expId.ExpressionId,
		Ttdo:         10,
		Status:       "Testing",
		Result:       2020,
	})
	require.NoError(t, err)

	exp, err = su.CalcClient.GetExpression(ctx, &calcv1.GetExpressionRequest{
		ExpressionId: expId.ExpressionId,
	})
	require.NoError(t, err)

	require.Equal(t, exp.Status, "Testing")
	require.Equal(t, exp.Result, int64(2020))
	require.Equal(t, exp.Ttdo, int64(10))

	_, err = su.CalcClient.UpdateArithmetic(ctx, &calcv1.UpdateArithmeticRequest{
		Sign:   "*",
		Ttdo:   25,
		UserId: int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	mult, err = su.CalcClient.GetArithmetic(ctx, &calcv1.GetArithmeticRequest{
		Sign:   "*",
		UserId: int64(claims["user_id"].(float64)),
	})
	require.NoError(t, err)

	require.Equal(t, mult.Ttdo, int64(25))
	require.Equal(t, "*", mult.Sign)
}

func TestCalc_Negative(t *testing.T) {
	ctx, su := suite.NewSuite(t)

	_, err := su.CalcClient.GetExpression(ctx, &calcv1.GetExpressionRequest{
		ExpressionId: 999,
	})
	require.Error(t, err)

	_, err = su.CalcClient.UpdateExpression(ctx, &calcv1.UpdateExpressionRequest{
		ExpressionId: 999,
	})
	require.Error(t, err)

	_, err = su.CalcClient.GetArithmetic(ctx, &calcv1.GetArithmeticRequest{
		Sign:   "*",
		UserId: 999,
	})
	require.Error(t, err)

	_, err = su.CalcClient.UpdateArithmetic(ctx, &calcv1.UpdateArithmeticRequest{
		Sign:   "*",
		UserId: 999,
	})
	require.Error(t, err)

}
