package web

import (
	"gitee.com/geekbang/basic-go/webook/internal/domain"
	"gitee.com/geekbang/basic-go/webook/internal/service"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	jwt "github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
	"unicode/utf8"
)

const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
	bizLogin             = "login"
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp
	svc            service.UserService
	codeSvc        service.CodeService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailRegexPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordRegexPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	// REST 风格
	//server.POST("/user", h.SignUp)
	//server.PUT("/user", h.SignUp)
	//server.GET("/users/:username", h.Profile)
	ug := server.Group("/users")
	// POST /users/signup
	ug.POST("/signup", h.SignUp)
	// POST /users/login
	//ug.POST("/login", h.Login)
	ug.POST("/login", h.LoginJWT)
	// POST /users/edit
	ug.POST("/edit", h.Edit)
	// GET /users/profile
	ug.GET("/profile", h.Profile)
	ug.POST("/login_sms/code/send", h.SendSMSLoginCode)
	ug.POST("/login_sms", h.LoginSMS)
}

func (h *UserHandler) SignUp(ctx *gin.Context) {
	type SignUpReq struct {
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}

	var req SignUpReq
	if err := ctx.Bind(&req); err != nil {
		return
	}

	isEmail, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isEmail {
		ctx.String(http.StatusOK, "非法邮箱格式")
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不对")
		return
	}

	isPassword, err := h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误")
		return
	}
	if !isPassword {
		ctx.String(http.StatusOK, "密码必须包含字母、数字、特殊字符，并且不少于八位")
		return
	}

	err = h.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})
	switch err {
	case nil:
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱冲突，请换一个")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) LoginJWT(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		uc := UserClaims{
			Uid:       u.Id,
			UserAgent: ctx.GetHeader("User-Agent"),
			RegisteredClaims: jwt.RegisteredClaims{
				// 30 分钟过期
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			},
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS512, uc)
		tokenStr, err := token.SignedString(JWTKey)
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
		}
		ctx.Header("x-jwt-token", tokenStr)
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	u, err := h.svc.Login(ctx, req.Email, req.Password)
	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", u.Id)
		sess.Options(sessions.Options{
			// 十分钟
			MaxAge: 30,
		})
		err = sess.Save()
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
			return
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或者密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
}

func (h *UserHandler) Edit(ctx *gin.Context) {
	type Request struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"aboutMe"`
	}
	var req Request
	if error := ctx.Bind(&req); error != nil {
		ctx.JSON(http.StatusOK, respJsonError(-1, "系统错误"))
		return
	}

	// 1. 昵称：字符串，你需要考虑允许的长度。
	if utf8.RuneCountInString(req.Nickname) == 0 {
		ctx.JSON(http.StatusOK, respJsonError(-1, "您的昵称不能输入空呢..."))
		return
	}
	if utf8.RuneCountInString(req.Nickname) > 128 {
		ctx.JSON(http.StatusOK, respJsonError(-1, "您的昵称不能大于128呢..."))
		return
	}
	// 2. 生日：前端输入为 1992-01-01 这种字符串。
	var (
		birthDay time.Time
		err      error
	)
	if birthDay, err = time.Parse(time.DateOnly, req.Birthday); err != nil {
		ctx.JSON(http.StatusOK, respJsonError(-1, "老哥，你输入到出生日期不符合 yyyy-MM-dd 格式呢..."))
		return
	}
	// 3. 个人简介：一段文本，你需要考虑允许的长度。
	if utf8.RuneCountInString(req.Nickname) > 500 {
		ctx.JSON(http.StatusOK, respJsonError(-1, "您的个人简介不能大于 500 个字符呢..."))
		return
	}

	h.withContextUser(ctx, func(uc *UserClaims) bool {
		err := h.svc.UpdateUserInfo(ctx, domain.User{
			Id:       uc.Uid,
			AboutMe:  req.AboutMe,
			Birthday: birthDay,
			Nickname: req.Nickname,
		})
		if err != nil {
			ctx.JSON(http.StatusOK, respJsonError(-1, "更新没有成功哦..."))
		} else {
			ctx.JSON(http.StatusOK, respJsonSuccess())
		}
		return true
	})
}

func respJsonError(code int64, msg string) gin.H {
	return gin.H{"code": -1, "msg": msg}
}

func respJsonSuccess() gin.H {
	return gin.H{"code": 0}
}

func (h *UserHandler) Profile(ctx *gin.Context) {
	if ok := h.withContextUser(ctx, func(uc *UserClaims) bool {
		if user, err := h.svc.FindById(ctx, uc.Uid); err == nil {
			ctx.JSON(http.StatusOK, user)
			return true
		}
		return false
	}); ok {
		return
	}
	ctx.String(http.StatusOK, "用户没有登录或者 session 过期了.")
}

func (h *UserHandler) withContextUser(ctx *gin.Context, callback func(*UserClaims) bool) bool {
	if ur, exists := ctx.Get("user"); exists == true {
		if uc, ok := ur.(UserClaims); ok == true {
			return callback(&uc)
		}
	}
	return false
}

func (h *UserHandler) SendSMSLoginCode(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		return
	}
	// 你这边可以校验 Req
	if req.Phone == "" {
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "请输入手机号码",
		})
		return
	}
	err := h.codeSvc.Send(ctx, bizLogin, req.Phone)
	switch err {
	case nil:
		ctx.JSON(http.StatusOK, Result{
			Msg: "发送成功",
		})
	case service.ErrCodeSendTooMany:
		// 事实上，防不住有人不知道怎么触发了
		// 少数这种错误，是可以接受的
		// 但是频繁出现，就代表有人在搞你的系统
		ctx.JSON(http.StatusOK, Result{
			Code: 4,
			Msg:  "短信发送太频繁，请稍后再试",
		})
	default:
		ctx.JSON(http.StatusOK, Result{
			Code: 5,
			Msg:  "系统错误",
		})
		// 补日志的
		log.Println(err)
	}

}

func (h *UserHandler) LoginSMS(context *gin.Context) {

}

var JWTKey = []byte("k6CswdUm77WKcbM68UQUuxVsHSpTCwgK")

type UserClaims struct {
	jwt.RegisteredClaims
	Uid       int64
	UserAgent string
}
