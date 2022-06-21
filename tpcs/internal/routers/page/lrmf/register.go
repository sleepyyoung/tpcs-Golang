package lrmf

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func RegisterHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "register.tmpl", gin.H{})
}
