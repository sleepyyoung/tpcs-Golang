package user

import (
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
	"time"
	"tpcs/global"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
	"tpcs/pkg/logger"
)

type User struct{}

func NewUser() User {
	return User{}
}

// UpdateUsername 修改用户名
func (u User) UpdateUsername(c *gin.Context) {
	response := app.NewResponse(c)
	newUsername, success := c.GetPostForm("new-username")
	if !success {
		response.ToFailResultResponse("修改用户名失败，表单异常，未获取新用户名！")
		return
	}
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.ToFailResultResponse("修改用户名失败，用户不存在！")
		return
	}
	userSvc := userService.New(c.Request.Context())
	user_, err := userSvc.GetUserById(userId)
	if err != nil {
		response.ToFailResultResponse("修改用户名失败，用户不存在！")
		return
	}
	user_.Username = &newUsername
	oldUser, _ := userSvc.GetUserByUsername(newUsername)
	if oldUser != nil {
		response.ToFailResultResponse("修改用户名失败，该用户名 " + newUsername + " 已存在！")
		return
	}
	modifyUser, err := userSvc.ModifyUser(user_)
	if err != nil {
		response.ToFailResultResponse("修改用户名失败！")
		return
	}
	if !modifyUser {
		response.ToFailResultResponse("修改用户名失败！")
		return
	} else {
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
			logger.Errorf("session.Save err: %v", err)
			return
		}
		response.ToSuccessResultResponse()
	}
}

// UpdateEmail 修改邮箱
func (u User) UpdateEmail(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())
	response := app.NewResponse(c)
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.ToFailResultResponse("修改邮箱失败，用户不存在！")
		return
	}
	newEmail, success := c.GetPostForm("email")
	if !success {
		response.ToFailResultResponse("修改用户名失败，表单异常，未获取新邮箱！")
		return
	}
	vcode, success := c.GetPostForm("vcode")
	if !success {
		response.ToFailResultResponse("修改用户名失败，表单异常，未获取验证码！")
		return
	}
	uu, _ := userSvc.GetUserByEmail(newEmail)
	if uu != nil {
		response.ToFailResultResponse("该邮箱已被其他用户绑定！")
		return
	}
	vv, _ := global.RedisClient.Get(newEmail).Result()
	if vcode != vv {
		response.ToFailResultResponse("邮箱验证码错误或已过期！")
		return
	}
	user_, err := userSvc.GetUserById(userId)
	if err != nil {
		response.ToFailResultResponse("修改邮箱失败，用户不存在！")
		return
	}
	user_.Email = &newEmail
	modifyUser, err := userSvc.ModifyUser(user_)
	if err != nil || !modifyUser {
		response.ToFailResultResponse("修改邮箱失败！")
		return
	}
	response.ToSuccessResultResponse()
}

// SendVcode 修改邮箱时发送验证码
func (u User) SendVcode(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())
	response := app.NewResponse(c)
	userId, err := strconv.Atoi(c.Param("userId"))
	if err != nil {
		response.ToFailResultResponse("发送验证码失败，用户不存在！")
		return
	}
	user_, err := userSvc.GetUserById(userId)
	if err != nil {
		response.ToFailResultResponse("发送验证码失败，用户不存在！")
		return
	}
	to, success := c.GetPostForm("email")
	if !success {
		response.ToFailResultResponse("表单异常，未获取到邮箱！")
		return
	}
	uu, _ := userSvc.GetUserByEmail(to)
	if uu != nil {
		response.ToFailResultResponse("该邮箱已被其他用户绑定！")
		return
	}

	vcode := fmt.Sprintf("%04v", rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(10000))
	_, err = global.RedisClient.Set(to, vcode, 1000000000*300).Result()
	if err != nil {
		response.ToFailResultResponse("验证码发送失败！")
		return
	}

	err = global.Email.SendMail(
		[]string{to},
		fmt.Sprintf("TPCS Verify Code for Modify Email"),
		fmt.Sprintf("您正在将该邮箱更新为TPCS用户 %v 的默认邮箱，验证码为 <strong>%v</strong> ，五分钟内有效", *user_.Username, vcode),
	)
	if err != nil {
		response.ToFailResultResponse("验证码发送失败！")
		return
	}
	response.ToSuccessResultResponse()
}
