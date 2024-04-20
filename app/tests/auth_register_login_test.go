package tests

import (
	"github.com/Smile8MrBread/gRPC_abstract_calculator/app/tests/suite"
	ssov1 "github.com/Smile8MrBread/gRPC_abstract_calculator/protos/gen/go/auth"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	appID      = 1
	emptyAppID = 0
	appSecret  = "verySecretApp"

	passDefLen = 10
)

func TestRegisterLogin_Happy(t *testing.T) {
	ctx, su := suite.NewSuite(t)

	email := gofakeit.Email()

	password := randomPassword()

	respReg, err := su.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := su.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Login:    email,
		Password: password,
		AppId:    appID,
	})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	assert.True(t, ok)

	assert.Equal(t, email, claims["email"].(string))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, respReg.GetUserId(), int64(claims["user_id"].(float64)))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(su.Cfg.TokenTTL).Unix(), claims["exp"].(float64), deltaSeconds)
}

func TestRegister_Dublicated(t *testing.T) {
	ctx, su := suite.NewSuite(t)

	email := gofakeit.Email()

	password := randomPassword()

	respReg, err := su.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respRegTwo, err := su.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Login:    email,
		Password: password,
	})
	require.Error(t, err)
	assert.Empty(t, respRegTwo.GetUserId())
	assert.ErrorContains(t, err, "user already exists")

}

func TestRegister_FailTests(t *testing.T) {
	ctx, su := suite.NewSuite(t)

	tests := []struct {
		name     string
		email    string
		password string
		expected string
	}{
		{
			name:     "empty email",
			email:    "",
			password: randomPassword(),
			expected: "email is required",
		},
		{
			name:     "empty password",
			email:    gofakeit.Email(),
			password: "",
			expected: "password is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := su.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Login:    tt.email,
				Password: tt.password,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expected)
		})
	}
}

func TestLogin_FailTests(t *testing.T) {
	ctx, su := suite.NewSuite(t)

	tests := []struct {
		name     string
		email    string
		password string
		appId    int32
		expected string
	}{
		{
			name:     "empty email",
			email:    "",
			password: randomPassword(),
			appId:    appID,
			expected: "email is required",
		},
		{
			name:     "empty password",
			email:    gofakeit.Email(),
			password: "",
			appId:    appID,
			expected: "password is required",
		},
		{
			name:     "empty appId",
			email:    gofakeit.Email(),
			password: randomPassword(),
			appId:    emptyAppID,
			expected: "appId is required",
		},
		{
			name:     "Non-matching password",
			email:    gofakeit.Email(),
			password: randomPassword(),
			appId:    appID,
			expected: "invalid email or password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := su.AuthClient.Register(ctx, &ssov1.RegisterRequest{
				Login:    gofakeit.Email(),
				Password: randomPassword(),
			})
			require.NoError(t, err)

			_, err = su.AuthClient.Login(ctx, &ssov1.LoginRequest{
				Login:    tt.email,
				Password: tt.password,
				AppId:    tt.appId,
			})
			require.Error(t, err)
			require.Contains(t, err.Error(), tt.expected)
		})
	}
}

func randomPassword() string {
	return gofakeit.Password(true, true, true, true, true, passDefLen)
}
