package mpage

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tpcs/global"
	"tpcs/internal/pojo/model"
	userService "tpcs/internal/service/user"
)

func IndexHandler(c *gin.Context) {
	userSvc := userService.New(c.Request.Context())

	ginH := gin.H{}
	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	user := userI.(model.User)
	ginH["username"] = user.Username
	isAdminUser, err := userSvc.IsAdminByUsernameAndPassword(*user.Username, *user.Password)
	if err != nil {
		global.Logger.Errorf("svc.IsAdminUser err: %v", err)
		return
	}
	if isAdminUser {
		ginH["isAdmin"] = true
	}

	c.HTML(http.StatusOK, "index.tmpl", ginH)
}
