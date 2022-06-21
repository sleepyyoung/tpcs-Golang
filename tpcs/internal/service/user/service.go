package user

import (
	"context"
	"fmt"
	"tpcs/global"
	"tpcs/internal/dao/user"
)

type Service struct {
	ctx context.Context
	dao *user.Dao
}

func New(ctx context.Context) Service {
	svc := Service{ctx: ctx}
	svc.dao = user.New(global.DBEngine)
	return svc
}

type BaseUserRequest struct {
	Username *string `form:"username" binding:"required"`
	Password *string `form:"password" binding:"required"`
}

type UserEmailRequest struct {
	Email *string `form:"email" binding:"required"`
	Vcode *string `form:"vcode" binding:"required"`
}

// UserLoginRequest 用户登录的请求体
type UserLoginRequest struct {
	// 用户名
	// 密码
	*BaseUserRequest
	// 是否保持30天内自动登录（on / off）
	StayLogin *string `form:"stay-login"`
}

// RegisterUserRequest 注册用户的请求体
type RegisterUserRequest struct {
	// 用户名
	// 密码
	*BaseUserRequest
	// 确认密码
	Password2 *string `form:"password2"`
	// 邮箱
	// 验证码
	*UserEmailRequest
	// 备注，管理员审核用
	Note *string `form:"note"`
}

func (rr *RegisterUserRequest) String() string {
	return fmt.Sprintf("{Username: %v, Password: %v, Password2: %v, Email: %v, Vcode: %v, Note: %v}", *rr.Username, *rr.Password, *rr.Password2, *rr.Email, *rr.Vcode, *rr.Note)
}

// ModifyPasswordRequest 修改密码的请求体
type ModifyPasswordRequest struct {
	// 用户名
	// 旧密码
	*BaseUserRequest
	// 新密码
	Password1 *string `json:"password1" binding:"required"`
	// 确认新密码
	Password2 *string `json:"password2" binding:"required"`
}

// ForgotPasswordRequest 忘记密码的请求体
type ForgotPasswordRequest struct {
	// 用户名
	// 新密码
	*BaseUserRequest
	// 确认新密码
	Password2 *string `json:"password2" binding:"required"`
	// 邮箱
	// 验证码
	*UserEmailRequest
}

// ThawUserRequest 解冻用户的请求体
type ThawUserRequest struct {
	Userid   *int    `form:"userid" binding:"required"`
	Username *string `form:"username" binding:"required"`
	Email    *string `form:"email" binding:"required"`
}

// FreezeUserRequest 冻结用户的请求体
type FreezeUserRequest struct {
	*ThawUserRequest
	Reason *string `form:"reason" binding:"required"`
}

// AuditUserRequest 审核用户的请求体
type AuditUserRequest struct {
	*ThawUserRequest
	Pass   *bool   `form:"pass" binding:"required"`
	Reason *string `form:"reason"`
}
