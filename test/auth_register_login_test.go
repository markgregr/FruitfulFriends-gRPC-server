package test

import (
	"github.com/brianvoe/gofakeit"
	"github.com/golang-jwt/jwt"
	"github.com/markgregr/FruitfulFriends-gRPC-server/test/suite"
	ssov1 "github.com/markgregr/FruitfulFriends-protos/gen/go/sso"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

const (
	emptyAppID     = 0
	appID          = 1
	appSecret      = "test-secret"
	passDefaultLen = 10
)

func TestRegisterLogin_Login_HappyPath(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	respReg, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)
	assert.NotEmpty(t, respReg.GetUserId())

	respLogin, err := st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    appID})
	require.NoError(t, err)

	loginTime := time.Now()

	token := respLogin.GetToken()
	require.NotEmpty(t, token)

	tokenParsed, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(appSecret), nil
	})
	require.NoError(t, err)

	claims, ok := tokenParsed.Claims.(jwt.MapClaims)
	require.True(t, ok)

	assert.Equal(t, respReg.GetUserId(), int64(claims["uid"].(float64)))
	assert.Equal(t, appID, int(claims["app_id"].(float64)))
	assert.Equal(t, email, claims["email"].(string))

	const deltaSeconds = 1

	assert.InDelta(t, loginTime.Add(st.Cfg.JWT.TokenTTL).Unix(), int64(claims["exp"].(float64)), deltaSeconds)

}

func TestRegisterLogin_Login_WrongPassword(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: "wrong-password",
		AppId:    appID})
	require.Error(t, err)
}

func TestRegisterLogin_Login_WrongEmail(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    "wrong-email",
		Password: password,
		AppId:    appID})
	require.Error(t, err)
}

func TestRegisterLogin_Login_EmptyEmail(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    "",
		Password: password,
		AppId:    appID})
	require.Error(t, err)
}

func TestRegisterLogin_Login_EmptyPassword(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: "",
		AppId:    appID})
	require.Error(t, err)
}

func TestRegisterLogin_DuplicatedRegistration(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.Error(t, err)
}

func TestRegisterLogin_Login_WrongAppID(t *testing.T) {
	ctx, st := suite.New(t)

	email := gofakeit.Email()
	password := randomFakePassword()

	_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
		Email:    email,
		Password: password,
	})
	require.NoError(t, err)

	_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
		Email:    email,
		Password: password,
		AppId:    emptyAppID})
	require.Error(t, err)
}

//func TestRegister_FailCases(t *testing.T) {
//	ctx, st := suite.New(t)
//
//	tests := []struct {
//		name     string
//		email    string
//		password string
//	}{
//		{
//			name:     "Register with Empty Password",
//			email:    gofakeit.Email(),
//			password: "",
//		},
//		{
//			name:     "Register with Empty Email",
//			email:    "",
//			password: randomFakePassword(),
//		},
//		{
//			name:     "Register with Both Empty",
//			email:    "",
//			password: "",
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
//				Email:    tt.email,
//				Password: tt.password,
//			})
//			require.Error(t, err)
//		})
//	}
//}
//
//func TestLogin_FailCases(t *testing.T) {
//	ctx, st := suite.New(t)
//
//	tests := []struct {
//		name     string
//		email    string
//		password string
//		appID    int32
//	}{
//		{
//			name:     "Login with Empty Password",
//			email:    gofakeit.Email(),
//			password: "",
//			appID:    appID,
//		},
//		{
//			name:     "Login with Empty Email",
//			email:    "",
//			password: randomFakePassword(),
//			appID:    appID,
//		},
//		{
//			name:     "Login with Both Empty Email and Password",
//			email:    "",
//			password: "",
//			appID:    appID,
//		},
//		{
//			name:     "Login with Non-Matching Password",
//			email:    gofakeit.Email(),
//			password: randomFakePassword(),
//			appID:    appID,
//		},
//		{
//			name:     "Login without AppID",
//			email:    gofakeit.Email(),
//			password: randomFakePassword(),
//			appID:    emptyAppID,
//		},
//	}
//
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			_, err := st.AuthClient.Register(ctx, &ssov1.RegisterRequest{
//				Email:    gofakeit.Email(),
//				Password: randomFakePassword(),
//			})
//			require.NoError(t, err)
//
//			_, err = st.AuthClient.Login(ctx, &ssov1.LoginRequest{
//				Email:    tt.email,
//				Password: tt.password,
//				AppId:    tt.appID,
//			})
//			require.Error(t, err)
//		})
//	}
//}

func randomFakePassword() string {
	return gofakeit.Password(true, true, true, true, false, passDefaultLen)
}
