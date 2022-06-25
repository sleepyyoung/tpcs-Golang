package user

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"tpcs/global"
	"tpcs/internal/pojo/model"
	"tpcs/internal/service"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
	"tpcs/pkg/email"
)

type Teacher struct{}

func NewTeacher() Teacher {
	return Teacher{}
}

// List 教师列表
func (t Teacher) List(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())

	var listRequest service.ListRequest
	err := c.ShouldBindQuery(&listRequest)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindQuery err: %v", err)
		return
	}

	var count int
	var teacherList []model.User
	teacherList, count, err = userSvc.NotAdminUserList(&listRequest)
	if err != nil {
		global.Logger.Errorf("svc.NotAdminUserList err: %v", err)
		return
	}

	var teacherMapList []map[string]interface{}
	for _, user_ := range teacherList {
		teacherMapList = append(teacherMapList, map[string]interface{}{
			"id":     *user_.Id,
			"name":   *user_.Username,
			"email":  *user_.Email,
			"status": *user_.Status,
		})
	}

	response := app.NewResponse(c)
	response.ToResponse(map[string]interface{}{
		"code":  0,
		"msg":   "",
		"count": count,
		"data":  teacherMapList,
	})
}

// Freeze 冻结用户
func (t Teacher) Freeze(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())
	response := app.NewResponse(c)

	var request userService.FreezeUserRequest
	err := c.ShouldBindWith(&request, binding.Form)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse("请稍后重试！")
		return
	}

	status := 2
	modifyUserSuccess, err := userSvc.ModifyUser(&model.User{Id: request.Userid, Status: &status})
	if err != nil || !modifyUserSuccess {
		global.Logger.Errorf("修改用户状态（待审核 -> 已冻结）失败，原因: %v", err)
		response.ToFailResultResponse("修改用户状态（待审核 -> 已冻结）失败！")
		return
	}

	err = email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	}).SendMail(
		[]string{*request.Email},
		fmt.Sprintf("TPCS 用户冻结通知"),
		fmt.Sprintf("您的TPCS账号 %v 已被冻结，解冻请联系管理员！原因：%v", *request.Username, *request.Reason),
	)
	if err != nil {
		global.Logger.Errorf("用户已冻结的邮件发送失败，原因: %v", err)
		response.ToFailResultResponse("用户已冻结，修改用户状态（待审核 -> 已冻结）成功，但是通知邮件发送失败！")
		return
	}
	response.ToSuccessResultResponse()
}

// Thaw 解冻用户
func (t Teacher) Thaw(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())
	response := app.NewResponse(c)

	var request userService.ThawUserRequest
	err := c.ShouldBindWith(&request, binding.Form)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse("请稍后重试！")
		return
	}

	status := 0
	modifyUserSuccess, err := userSvc.ModifyUser(&model.User{Id: request.Userid, Status: &status})
	if err != nil || !modifyUserSuccess {
		global.Logger.Errorf("修改用户状态（已冻结 -> 正常）失败，原因: %v", err)
		response.ToFailResultResponse("修改用户状态（已冻结 -> 正常）失败！")
		return
	}

	err = email.NewEmail(&email.SMTPInfo{
		Host:     global.EmailSetting.Host,
		Port:     global.EmailSetting.Port,
		IsSSL:    global.EmailSetting.IsSSL,
		UserName: global.EmailSetting.UserName,
		Password: global.EmailSetting.Password,
		From:     global.EmailSetting.From,
	}).SendMail(
		[]string{*request.Email},
		fmt.Sprintf("TPCS 用户解冻通知"),
		fmt.Sprintf("您的TPCS账号 %v 已解冻！", *request.Username),
	)
	if err != nil {
		global.Logger.Errorf("用户已解冻的邮件发送失败，原因: %v", err)
		response.ToFailResultResponse("用户已解冻，修改用户状态（已冻结 -> 正常）成功，但是通知邮件发送失败！")
		return
	}
	response.ToSuccessResultResponse()
}

func (t Teacher) Audit2(c *gin.Context) {
	t.Audit(c)
}

// Audit 审核用户
func (t Teacher) Audit(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())
	response := app.NewResponse(c)

	var request userService.AuditUserRequest
	err := c.ShouldBindWith(&request, binding.Form)
	if err != nil {
		global.Logger.Errorf("c.ShouldBindWith err: %v", err)
		response.ToFailResultResponse("请稍后重试！")
		return
	}

	if *request.Pass {
		status := 0
		modifyUserSuccess, err := userSvc.ModifyUser(&model.User{Id: request.Userid, Status: &status})
		if err != nil {
			global.Logger.Errorf("修改用户状态（待审核 -> 正常）失败，原因: %v", err)
			response.ToFailResultResponse("修改用户状态（待审核 -> 正常）失败！")
			return
		}
		if modifyUserSuccess {
			err = email.NewEmail(&email.SMTPInfo{
				Host:     global.EmailSetting.Host,
				Port:     global.EmailSetting.Port,
				IsSSL:    global.EmailSetting.IsSSL,
				UserName: global.EmailSetting.UserName,
				Password: global.EmailSetting.Password,
				From:     global.EmailSetting.From,
			}).SendMail(
				[]string{*request.Email},
				fmt.Sprintf("TPCS 用户注册审核通过通知"),
				fmt.Sprintf("您近期注册的TPCS用户 %v 已通过审核，欢迎使用！", *request.Username),
			)
			if err != nil {
				global.Logger.Errorf("用户通过审核的邮件发送失败，原因: %v", err)
				response.ToFailResultResponse("用户已通过审核，修改用户状态（待审核 -> 正常）成功，但是通知邮件发送失败！")
				return
			}
		} else {
			response.ToFailResultResponse("修改用户状态（待审核 -> 正常）失败！")
			return
		}
	} else {
		isDelete, err := userSvc.DeleteUserById(*request.Userid)
		if err != nil {
			global.Logger.Errorf("删除未通过审核的用户信息失败，原因: %v", err)
			response.ToFailResultResponse("删除未通过审核的用户信息失败！")
			return
		}
		if isDelete {
			err = email.NewEmail(&email.SMTPInfo{
				Host:     global.EmailSetting.Host,
				Port:     global.EmailSetting.Port,
				IsSSL:    global.EmailSetting.IsSSL,
				UserName: global.EmailSetting.UserName,
				Password: global.EmailSetting.Password,
				From:     global.EmailSetting.From,
			}).SendMail(
				[]string{*request.Email},
				fmt.Sprintf("TPCS 用户注册审核未通过"),
				fmt.Sprintf("您近期注册的TPCS用户 %v 审核未通过，请重新注册，原因：<textarea>%v</textarea>", *request.Username, *request.Reason),
			)
			if err != nil {
				global.Logger.Errorf("用户通过审核的邮件发送失败，原因: %v", err)
				response.ToFailResultResponse("用户已通过审核，修改用户状态（待审核 -> 正常）成功，但是通知邮件发送失败！")
				return
			}
		} else {
			response.ToFailResultResponse("删除未通过审核的用户信息失败！")
			return
		}
	}
	response.ToSuccessResultResponse()
}
