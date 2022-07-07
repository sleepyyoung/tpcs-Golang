package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/app"
	"tpcs/pkg/logger"
)

func TeacherListHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "teacher-list.tmpl", gin.H{})
}

func TeacherAuditBySessionHandler(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		logger.Errorf("strconv.Atoi(c.Param(\"id\")) error: %v\n", err)
		return
	}

	user, err := userSvc.GetUserById(id)
	if err != nil {
		logger.Errorf("userSvc.GetUserById error: %v\n", err)
		return
	}
	c.HTML(http.StatusOK, "teacher-audit.tmpl", gin.H{"user": user})
}

func TeacherAuditByJWTHandler(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())

	calaims, err := app.ParseToken(c.Query("token"))
	if err != nil {
		c.HTML(http.StatusOK, "404.tmpl", gin.H{})
		logger.Errorf("app.ParseToken error: %v\n", err)
		return
	}
	if calaims == nil {
		c.HTML(http.StatusOK, "404.tmpl", gin.H{})
		return
	}

	user, err := userSvc.GetUserByUsername(calaims.Audience)
	if err != nil {
		logger.Errorf("userSvc.GetUserByUsername error: %v\n", err)
		c.HTML(http.StatusOK, "404.tmpl", gin.H{"message": "该用户已被其他管理员审核，审核结果：未通过"})
		return
	}
	if *user.Status == 0 {
		c.HTML(http.StatusOK, "404.tmpl", gin.H{"message": "该用户已被其他管理员审核，审核结果：通过"})
		return
	}

	c.HTML(http.StatusOK, "teacher-audit2.tmpl", gin.H{"user": user})
}
