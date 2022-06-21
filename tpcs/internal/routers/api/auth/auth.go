package auth

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"math/rand"
	"net/http"
	"strings"
	"time"
	"tpcs/global"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service/user"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
)

type Auth struct{}

func NewAuth() Auth {
	return Auth{}
}

// SendVcode4Register 发送验证码（注册用）
func (a Auth) SendVcode4Register(c *gin.Context) {
	response := app.NewResponse(c)
	email_, get := c.GetPostForm("email")
	if !get {
		response.ToFailResultResponse(pojo.ResultMsg_FormErr_NoneEmail)
		return
	}
	vcode := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	_, err := global.RedisClient.Set(email_, vcode, 1000000000*300).Result()
	if err != nil {
		global.Logger.Errorf("global.RedisClient.Set err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_SendVcodeFail)
		return
	}

	err = global.Email.SendMail(
		[]string{email_},
		fmt.Sprintf("TPCS Verify Code for Register"),
		fmt.Sprintf("您正在注册为TPCS用户，验证码为 <strong>%v</strong> ，五分钟内有效", vcode),
	)
	if err != nil {
		global.Logger.Errorf("SendMail err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_SendVcodeFail)
		return
	}
	response.ToSuccessResultResponse()
}

// SendVcode4Forgot 发送验证码（忘记密码用）
func (a Auth) SendVcode4Forgot(c *gin.Context) {
	response := app.NewResponse(c)
	to, get := c.GetPostForm("email")
	if !get {
		response.ToFailResultResponse(pojo.ResultMsg_FormErr_NoneEmail)
		return
	}
	username, get := c.GetPostForm("username")
	if !get {
		response.ToFailResultResponse(pojo.ResultMsg_FormErr_NoneUsername)
		return
	}
	vcode := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	_, err := global.RedisClient.Set(to, vcode, 1000000000*300).Result()
	if err != nil {
		global.Logger.Errorf("global.RedisClient.Set err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_SendVcodeFail)
		return
	}

	err = global.Email.SendMail(
		[]string{to},
		fmt.Sprintf("TPCS Verify Code for Modify Password"),
		fmt.Sprintf("TPCS用户 %v 正在修改密码，验证码为 <strong>%v</strong> ，五分钟内有效", username, vcode),
	)
	if err != nil {
		global.Logger.Errorf("SendMail err: %v\n", err)
		response.ToFailResultResponse(pojo.ResultMsg_SendVcodeFail)
		return
	}
	response.ToSuccessResultResponse()
}

// Register 注册
func (a Auth) Register(c *gin.Context) {
	response := app.NewResponse(c)
	var registerUserRequest *user.RegisterUserRequest
	err := c.ShouldBindWith(&registerUserRequest, binding.Form)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}
	if registerUserRequest.Username == nil || strings.Trim(*registerUserRequest.Username, " ") == "" {
		response.ToFailResultResponse(pojo.ResultMsg_UsernameNotNone)
		return
	}
	if len(strings.Trim(*registerUserRequest.Username, " ")) > 15 {
		response.ToFailResultResponse(pojo.ResultMsg_UsernameLengthLessThan15)
		return
	}
	if registerUserRequest.Email == nil || strings.Trim(*registerUserRequest.Email, " ") == "" {
		response.ToFailResultResponse(pojo.ResultMsg_EmailNotNone)
		return
	}
	if registerUserRequest.Password == nil || strings.Trim(*registerUserRequest.Password, " ") == "" {
		response.ToFailResultResponse(pojo.ResultMsg_PasswordNotNone)
		return
	}
	if registerUserRequest.Password2 == nil || strings.Trim(*registerUserRequest.Password2, " ") == "" {
		response.ToFailResultResponse(pojo.ResultMsg_Password2NotNone)
		return
	}
	if strings.Trim(*registerUserRequest.Password, " ") != strings.Trim(*registerUserRequest.Password2, " ") {
		response.ToFailResultResponse(pojo.ResultMsg_2PasswordNotSame)
		return
	}

	userSvc := userService.New(c.Request.Context())
	u, _ := userSvc.GetUserByUsernameAndPassword(*registerUserRequest.BaseUserRequest.Username, *registerUserRequest.BaseUserRequest.Password)
	if u != nil {
		//response.ToFailResultResponse("用户 " + *registerUserRequest.Username + " 已存在！")
		response.ToFailResultResponse(pojo.ResultMsg_UserExisted)
		return
	}
	if registerUserRequest.Email == nil {
		response.ToFailResultResponse(pojo.ResultMsg_EmailNotNone)
		return
	}
	vcode, _ := global.RedisClient.Get(*registerUserRequest.Email).Result()
	if registerUserRequest.Vcode == nil {
		response.ToFailResultResponse(pojo.ResultMsg_EmailVcodeNotNone)
		return
	}
	if vcode != *registerUserRequest.Vcode {
		response.ToFailResultResponse(pojo.ResultMsg_EmailVcodeIllegal)
		return
	}
	u, _ = userSvc.GetUserByEmail(*registerUserRequest.Email)
	if u != nil {
		//response.ToFailResultResponse("该邮箱 " + *registerUserRequest.Email + " 已被其他用户绑定！")
		response.ToFailResultResponse(pojo.ResultMsg_EmailExisted)
		return
	}

	status := 1
	isAdministrator := false
	newUser := &model.User{
		Username:        registerUserRequest.Username,
		Password:        registerUserRequest.Password,
		Email:           registerUserRequest.Email,
		Note:            registerUserRequest.Note,
		Status:          &(status),
		IsAdministrator: &(isAdministrator),
	}

	err = userSvc.CreateUser(newUser)
	if err != nil {
		global.Logger.Errorf("svc.CreateUser err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
		return
	} else {
		err := global.RedisClient.Publish(global.RedisSetting.Topic, newUser).Err()
		if err != nil {
			global.Logger.Errorf("Redis Publish error: %v", err)
		}
	}
	response.ToSuccessResultResponse()
}

// Login 登录
func (a Auth) Login(c *gin.Context) {
	var userRequest *user.UserLoginRequest
	err := c.ShouldBindWith(&userRequest, binding.Form)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title": "表单异常",
		})
		return
	}
	userSvc := userService.New(c.Request.Context())
	u, _ := userSvc.GetUserByUsernameAndPassword(*userRequest.BaseUserRequest.Username, *userRequest.BaseUserRequest.Password)
	if u == nil {
		c.HTML(http.StatusOK, "login.tmpl", gin.H{
			"title":    "用户名或密码输入错误",
			"username": *userRequest.Username,
		})
	} else {
		session := sessions.DefaultMany(c, "user")
		session.Set("user", u)
		var maxAge int
		if *userRequest.StayLogin == "on" {
			maxAge = 30 * 24 * 60 * 60
		} else {
			maxAge = 0
		}
		options := sessions.Options{
			Path:     "/",
			Domain:   "",
			MaxAge:   maxAge,
			Secure:   false,
			HttpOnly: true,
			SameSite: http.SameSiteDefaultMode,
		}
		session.Options(options)
		err := session.Save()
		if err != nil {
			global.Logger.Errorf("session.Save err: %v", err)
			return
		}
		c.Redirect(http.StatusFound, "/")
	}
}

// ForgotPassword 忘记密码
func (a Auth) ForgotPassword(c *gin.Context) {
	response := app.NewResponse(c)
	var forgotPasswordRequest *user.ForgotPasswordRequest
	err := c.ShouldBindJSON(&forgotPasswordRequest)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	vcode, _ := global.RedisClient.Get(*forgotPasswordRequest.Email).Result()
	if vcode != *forgotPasswordRequest.Vcode {
		response.ToFailResultResponse(pojo.ResultMsg_EmailVcodeIllegal)
		return
	}

	userSvc := userService.New(c.Request.Context())
	u, _ := userSvc.GetUserByUsername(*forgotPasswordRequest.Username)
	if u == nil {
		//response.ToFailResultResponse("修改密码失败，用户 " + *forgotPasswordRequest.Username + " 不存在！")
		response.ToFailResultResponse(pojo.ResultMsg_UserNotExists)
		return
	} else {
		u.Password = forgotPasswordRequest.Password
		modify, _ := userSvc.ModifyUser(u)
		if modify {
			session := sessions.DefaultMany(c, "user")
			session.Delete("user")
			options := sessions.Options{
				Path:     "/",
				Domain:   "",
				MaxAge:   -1,
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteDefaultMode,
			}
			session.Options(options)
			err := session.Save()
			if err != nil {
				global.Logger.Errorf("session.Save err: %v", err)
				return
			}
		} else {
			response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
			return
		}
	}
	response.ToSuccessResultResponse()
}

// ModifyPassword 修改密码
func (a Auth) ModifyPassword(c *gin.Context) {
	response := app.NewResponse(c)
	var modifyPasswordRequest *user.ModifyPasswordRequest
	err := c.ShouldBindJSON(&modifyPasswordRequest)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse(pojo.ResultMsg_FormParseErr)
		return
	}

	userSvc := userService.New(c.Request.Context())
	u, _ := userSvc.GetUserByUsernameAndPassword(*modifyPasswordRequest.BaseUserRequest.Username, *modifyPasswordRequest.BaseUserRequest.Password)
	if u == nil {
		response.ToFailResultResponse(pojo.ResultMsg_UsernameOrPasswordErr)
		return
	} else {
		u.Password = modifyPasswordRequest.Password1
		modify, _ := userSvc.ModifyUser(u)
		if modify {
			session := sessions.DefaultMany(c, "user")
			session.Delete("user")
			options := sessions.Options{
				Path:     "/",
				Domain:   "",
				MaxAge:   -1,
				Secure:   false,
				HttpOnly: true,
				SameSite: http.SameSiteDefaultMode,
			}
			session.Options(options)
			err := session.Save()
			if err != nil {
				global.Logger.Errorf("session.Save err: %v", err)
				return
			}
		} else {
			response.ToFailResultResponse(pojo.ResultMsg_TryAgainLater)
			return
		}
	}
	response.ToSuccessResultResponse()
}

func (a Auth) Logout(c *gin.Context) {
	response := app.NewResponse(c)
	session := sessions.DefaultMany(c, "user")
	session.Delete("user")
	options := sessions.Options{
		Path:     "/",
		Domain:   "",
		MaxAge:   -1,
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteDefaultMode,
	}
	session.Options(options)
	err := session.Save()
	if err != nil {
		global.Logger.Errorf("session.Save err: %v", err)
		return
	}

	response.ToSuccessResultResponse()
}
