package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"tpcs/internal/pojo"
	"tpcs/internal/pojo/model"
	userService "tpcs/internal/service/user"
)

func UserPageInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSvc := userService.New(c.Request.Context())

		session := sessions.DefaultMany(c, "user")
		userI := session.Get("user")
		if userI == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			user, err := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
			if user == nil || err != nil {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
			}
			if *user.Status == 2 {
				c.HTML(http.StatusOK, "login.tmpl", gin.H{
					"title": pojo.ResultMsg_UserFreezed,
				})
				c.Abort()
			}
		}
		c.Next()
	}
}

func TeacherPageInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSvc := userService.New(c.Request.Context())

		session := sessions.DefaultMany(c, "user")
		userI := session.Get("user")
		if userI == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			user, err := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
			if user == nil || err != nil {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
			}
			if *user.IsAdministrator {
				c.Redirect(http.StatusFound, "/403")
				c.Abort()
			}
			if *user.Status == 2 {
				c.HTML(http.StatusOK, "login.tmpl", gin.H{
					"title": pojo.ResultMsg_UserFreezed,
				})
				c.Abort()
			}
		}
		c.Next()
	}
}

func AdminPageInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSvc := userService.New(c.Request.Context())

		session := sessions.DefaultMany(c, "user")
		userI := session.Get("user")
		if userI == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
		} else {
			user, err := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
			if user == nil || err != nil {
				c.Redirect(http.StatusFound, "/login")
				c.Abort()
			}
			if !*user.IsAdministrator {
				c.Redirect(http.StatusFound, "/403")
				c.Abort()
			}
			if *user.Status == 2 {
				c.HTML(http.StatusOK, "login.tmpl", gin.H{
					"title": pojo.ResultMsg_UserFreezed,
				})
				c.Abort()
			}
		}
		c.Next()
	}
}

func UserApiInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSvc := userService.New(c.Request.Context())

		session := sessions.DefaultMany(c, "user")
		userI := session.Get("user")
		if userI == nil {
			c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
			c.Abort()
		} else {
			user, err := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
			if user == nil || err != nil {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
				c.Abort()
			}
			if *user.Status == 2 {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_UserFreezed})
				c.Abort()
			}
		}
		c.Next()
	}
}

func TeacherApiInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSvc := userService.New(c.Request.Context())

		session := sessions.DefaultMany(c, "user")
		userI := session.Get("user")
		if userI == nil {
			c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
			c.Abort()
		} else {
			user, err := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
			if user == nil || err != nil {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
				c.Abort()
			}
			if *user.IsAdministrator {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
				c.Abort()
			}
			if *user.Status == 2 {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_UserFreezed})
				c.Abort()
			}
		}
		c.Next()
	}
}

func AdminApiInterceptor() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSvc := userService.New(c.Request.Context())

		session := sessions.DefaultMany(c, "user")
		userI := session.Get("user")
		if userI == nil {
			c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
			c.Abort()
		} else {
			user, err := userSvc.GetUserByUsernameAndPassword(*userI.(model.User).Username, *userI.(model.User).Password)
			if user == nil || err != nil {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
				c.Abort()
			}
			if !*user.IsAdministrator {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_InsufficientPermissions})
				c.Abort()
			}
			if *user.Status == 2 {
				c.JSON(http.StatusForbidden, pojo.Result{Success: pojo.ResultSuccess_False, Msg: pojo.ResultMsg_UserFreezed})
				c.Abort()
			}
		}
		c.Next()
	}
}
