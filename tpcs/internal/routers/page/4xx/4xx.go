package _4xx

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func NotFoundHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "404.tmpl", gin.H{})
}

func ForbiddenHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "403.tmpl", gin.H{})
}
