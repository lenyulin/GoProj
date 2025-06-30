package web

import (
	"GoProj/wedy/internal/domian"
	"GoProj/wedy/internal/service"
	"errors"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"
)

const (
	phoneRegexPattern    = `^1[3-9]\d{9}$`
	passwordRegexPattern = `^[a-zA-Z0-9]+$`
	bizLogin             = "login"
)

var (
	ErrDuplicatedUser        = service.ErrDuplicatedUser
	ErrInvalidUserOrPassword = service.ErrInvalidUserOrPassword
)

type UeerHandler struct {
	phoneRegex    *regexp.Regexp
	passwordRegex *regexp.Regexp
	svc           service.UserService
	codeService   service.CodeService
}

func NewUserHandler(svc service.UserService, codeService service.CodeService) *UeerHandler {
	return &UeerHandler{
		phoneRegex:    regexp.MustCompile(phoneRegexPattern, regexp.None),
		passwordRegex: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:           svc,
		codeService:   codeService,
	}
}

func (h *UeerHandler) RegisiterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", h.SignUp)
	ug.POST("/login", h.Login)
	ug.POST("/loginjwt", h.LoginJWT)
	ug.POST("/loginsms/code/send", h.SendSMSLoginCode)
	ug.POST("/loginsms/", h.LoginSMS)
}

func (h *UeerHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Phone           string `json:"Phone"`
		Password        string `json:"Password"`
		ConfirmPassword string `json:"ConfirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "Internal Server Error")
		return
	}
	phoneOK, err := h.phoneRegex.MatchString(req.Phone)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	if !phoneOK {
		ctx.String(http.StatusOK, "Phone number is not valid")
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "Password is not equal")
		return
	}
	isPassword, err := h.passwordRegex.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "Internal Server Error")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "Password is not valid")
		return
	}
	err = h.svc.Signup(ctx, domian.User{
		Phone:    req.Phone,
		Password: req.Password,
	})
	switch {
	case err == nil:
		ctx.String(http.StatusOK, "User Signed Up Succeed")
	case errors.Is(err, ErrDuplicatedUser):
		ctx.String(http.StatusOK, "User Exists")
	default:
		ctx.String(http.StatusOK, "Internal Server Error")
	}
}

func (h *UeerHandler) Login(ctx *gin.Context) {
	type SignUpReq struct {
		Phone    string `json:"Phone"`
		Password string `json:"Password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(400, "Internal Server Error")
		return
	}
	u, err := h.svc.Login(ctx, req.Phone, req.Password)
	switch {
	case err == nil:
		//session
		session := sessions.Default(ctx)
		session.Set("userId", u.Id)
		err = session.Save()
		if err != nil {
			ctx.String(http.StatusOK, "Internal Server Error")
			return
		}
		ctx.String(http.StatusOK, "User Login Succeed")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "Invalid User or Password")
	default:
		ctx.String(http.StatusOK, "Internal Server Error")
	}
	return
}

type UserClaims struct {
	jwt.RegisteredClaims
	UserId    int64
	UserAgent string
}

var SigningKey = []byte("7TUqGzmiBIj0Fcl8FbBd4wqN7xRx9SrP")

func (h *UeerHandler) LoginJWT(ctx *gin.Context) {
	type SignUpReq struct {
		Phone    string `json:"Phone"`
		Password string `json:"Password"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.String(http.StatusBadRequest, "Internal Server Error")
		return
	}
	u, err := h.svc.Login(ctx, req.Phone, req.Password)
	switch {
	case err == nil:
		//session
		h.setJWTToken(ctx, u.Id)
		ctx.String(http.StatusOK, "User Login Succeed")
	case errors.Is(err, service.ErrInvalidUserOrPassword):
		ctx.String(http.StatusOK, "Invalid User or Password")
	default:
		ctx.JSON(http.StatusOK, "Internal Server Error")
	}
	return
}

func (h *UeerHandler) LoginSMS(ctx *gin.Context) {
	type SignUpReq struct {
		Phone string `json:"Phone"`
		Code  string `json:"code"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, "Internal Server Error")
		return
	}
	ok, err := h.codeService.Verify(ctx, bizLogin, req.Phone, req.Code)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "Internal Server Error",
		})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Code Invalid",
		})
		return
	}
	u, err := h.svc.FindOrCreate(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "Internal Server Error",
		})
		return
	}
	h.setJWTToken(ctx, u.Id)
	ctx.JSON(http.StatusOK, Result{
		Code: 4,
		Msg:  "Login Success",
	})
}
func (h *UeerHandler) setJWTToken(ctx *gin.Context, uid int64) {
	uc := UserClaims{
		UserId:    uid,
		UserAgent: ctx.GetHeader("User-Agent"),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		}}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
	tokenStr, err := token.SignedString([]byte(SigningKey))
	if err != nil {
		ctx.JSON(http.StatusOK, "Internal Server Error")
	}
	ctx.Header("x-jwt-token", tokenStr)
}

func (h *UeerHandler) SendSMSLoginCode(ctx *gin.Context) {
	type SignUpReq struct {
		Phone string `json:"Phone"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		ctx.JSON(http.StatusOK, "Internal Server Error")
		return
	}
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Phone number is empty",
		})
		return
	}
	if regexRes, _ := h.phoneRegex.MatchString(req.Phone); regexRes != true {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "Phone number not invalid format",
		})
		return
	}
	err := h.codeService.Send(ctx, bizLogin, req.Phone)
	switch {
	case err == nil:
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "SMS Send Succeed",
		})
	case errors.Is(err, service.ErrCodeVerifySendTooMany):
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "SMS Send Too Many",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "Internal Server Error",
		})
	}
}
