package lrmf

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ModifyPasswordHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "modify-password.tmpl", gin.H{"username": c.Param("username")})
}
