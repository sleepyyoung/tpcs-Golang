package course

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func CourseListHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "course-list.tmpl", gin.H{})
}
