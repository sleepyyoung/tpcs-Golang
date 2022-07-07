package user

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tpcs/internal/pojo/model"
	userService "tpcs/internal/service/user"
	"tpcs/pkg/logger"
)

func UserHandler(c *gin.Context) {
	userSvc := userService.New(c)
	ginH := gin.H{
		"isMe": false,
	}

	userByUsername, err := userSvc.GetUserByUsername(c.Param("username"))
	if err != nil {
		logger.Errorf("userSvc.GetUserByUsername err: %v\n", err)
		return
	}

	session := sessions.DefaultMany(c, "user")
	userI := session.Get("user")
	userBySession := userI.(model.User)
	if *userByUsername.Username == *userBySession.Username {
		ginH["isMe"] = true
	}
	ginH["user"] = userByUsername
	c.HTML(http.StatusOK, "user-info.tmpl", ginH)
}
