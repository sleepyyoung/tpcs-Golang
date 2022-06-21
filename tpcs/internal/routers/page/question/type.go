package question

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// QuestionTypeListHandler 题型管理页面
func QuestionTypeListHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "question-type-list.tmpl", gin.H{})
}
