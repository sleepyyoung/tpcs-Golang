package lrmf

import (
	"github.com/gin-gonic/gin"
	"net/http"
	userService "tpcs/internal/service/user"
)

func ForgotPasswordHandler(c *gin.Context) {
	username := c.Param("username")
	userSvc := userService.New(c.Request.Context())

	user, _ := userSvc.GetUserByUsername(username)
	if user == nil {
		c.HTML(http.StatusOK, "404.tmpl", gin.H{})
	} else {
		c.HTML(http.StatusOK, "forgot-password.tmpl", gin.H{
			"username": *user.Username,
			"email":    *user.Email,
		})
	}
}
