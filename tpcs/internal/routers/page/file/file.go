package file

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tpcs/internal/pojo/model"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/logger"
)

// FileListHandler 文件列表
func FileListHandler(c *gin.Context) {
	ginH := gin.H{}
	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)
	userSvc := userService.New(c.Request.Context())

	ginH["username"] = user.Username
	isAdminUser, err := userSvc.IsAdminByUsernameAndPassword(*user.Username, *user.Password)
	if err != nil {
		logger.Errorf("svc.IsAdminUser err: %v", err)
		return
	}
	if isAdminUser {
		ginH["isAdmin"] = true
	}

	c.HTML(http.StatusOK, "file-list.tmpl", ginH)
}

// FileUploadHandler 文件上传
func FileUploadHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "file-upload.tmpl", gin.H{})
}
