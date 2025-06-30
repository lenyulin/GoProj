package web

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/service"
	svcmock "GoProj/wedy/internal/service/mocks"
	"bytes"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/assert/v2"
	assert2 "github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserSignUpHdl(t *testing.T) {
	testcases := []struct {
		name       string
		wantCode   int
		wantBody   string
		reqBuilder func(t *testing.T) *http.Request
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
	}{
		{
			name:     "Ok",
			wantCode: http.StatusOK,
			wantBody: "User Signed Up Succeed",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domian.User{
					Phone:    "13333333333",
					Password: "hello123",
				}).Return(nil)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/signup",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"hello123","ConfirmPassword":"hello123"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "ErrDuplicatedUser",
			wantCode: http.StatusOK,
			wantBody: "User Exists",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Signup(gomock.Any(), domian.User{
					Phone:    "13333333333",
					Password: "hello123",
				}).Return(ErrDuplicatedUser)
				return userSvc, nil
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/signup",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"hello123","ConfirmPassword":"hello123"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "PasswordNotValid",
			wantCode: http.StatusOK,
			wantBody: "Password is not valid",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/signup",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"。。。","ConfirmPassword":"。。。"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "Password is not equal",
			wantCode: http.StatusOK,
			wantBody: "Password is not equal",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/signup",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"1234","ConfirmPassword":"123"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "Internal Server Error",
			wantCode: http.StatusBadRequest,
			wantBody: "Internal Server Error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/signup",
					bytes.NewBuffer([]byte(`{"Phone":"1333333333","Password":"1234","ConfirmPassword":"123`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "Phone number is not valid",
			wantCode: http.StatusOK,
			wantBody: "Phone number is not valid",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/signup",
					bytes.NewBuffer([]byte(`{"Phone":"123456","Password":"1234","ConfirmPassword":"1234"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)
			server := gin.Default()
			hdl.RegisiterRoutes(server)
			req := tc.reqBuilder(t)
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserLoginSessionHdl(t *testing.T) {
	testcases := []struct {
		name       string
		wantCode   int
		wantBody   string
		reqBuilder func(t *testing.T) *http.Request
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
	}{
		{
			name:     "Invalid User or Password",
			wantCode: http.StatusOK,
			wantBody: "Invalid User or Password",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "13333333333", "hello123").Return(
					domian.User{}, ErrInvalidUserOrPassword)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/login",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"hello123"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "Login Success",
			wantCode: http.StatusOK,
			wantBody: "User Login Succeed",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "13333333333", "hello123").Return(
					domian.User{
						Phone: "13333333333",
					}, nil)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/login",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"hello123"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
		{
			name:     "Bind Error",
			wantCode: http.StatusBadRequest,
			wantBody: "Internal Server Error",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(http.MethodPost, "/user/login",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"hello1"`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)
			server := gin.Default()
			store := cookie.NewStore([]byte("your-secret-key"))
			server.Use(sessions.Sessions("userId", store))
			hdl.RegisiterRoutes(server)
			req := tc.reqBuilder(t)
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}

func TestUserLoginJWTHdl(t *testing.T) {
	testcases := []struct {
		name       string
		wantCode   int
		wantBody   string
		reqBuilder func(t *testing.T) *http.Request
		mock       func(ctrl *gomock.Controller) (service.UserService, service.CodeService)
	}{
		{
			name:     "User Login Succeed",
			wantCode: http.StatusOK,
			wantBody: "User Login Succeed",
			mock: func(ctrl *gomock.Controller) (service.UserService, service.CodeService) {
				userSvc := svcmock.NewMockUserService(ctrl)
				userSvc.EXPECT().Login(gomock.Any(), "13333333333", "hello123").Return(
					domian.User{Phone: "13333333333"}, nil)
				codeSvc := svcmock.NewMockCodeService(ctrl)
				return userSvc, codeSvc
			},
			reqBuilder: func(t *testing.T) *http.Request {
				req, err := http.NewRequest(
					http.MethodPost,
					"/user/loginjwt",
					bytes.NewBuffer([]byte(`{"Phone":"13333333333","Password":"hello123"}`)))
				req.Header.Add("Content-Type", "application/json")
				assert2.NoError(t, err)
				return req
			},
		},
	}
	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			userSvc, codeSvc := tc.mock(ctrl)
			hdl := NewUserHandler(userSvc, codeSvc)

			server := gin.Default()
			//server.Use((&middlewares.LoginJWTMiddleWareBuilder{}).CheckJWTLogin())
			hdl.RegisiterRoutes(server)

			req := tc.reqBuilder(t)
			resp := httptest.NewRecorder()
			server.ServeHTTP(resp, req)
			assert.Equal(t, tc.wantCode, resp.Code)
			assert.Equal(t, tc.wantBody, resp.Body.String())
		})
	}
}
