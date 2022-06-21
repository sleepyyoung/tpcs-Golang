package user

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	userService "tpcs/internal/service/user"
)

func TeacherListHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "teacher-list.tmpl", gin.H{})
}

func TeacherAuditHandler(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return
	}

	user, err := userSvc.GetUserById(id)
	if err != nil {
		return
	}
	c.HTML(http.StatusOK, "teacher-audit.tmpl", gin.H{"user": user})
}
