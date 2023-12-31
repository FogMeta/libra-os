package router

import (
	"net/http"
	"strings"

	"github.com/FogMeta/libra-os/api/result"
	apiV1 "github.com/FogMeta/libra-os/api/v1"
	"github.com/FogMeta/libra-os/service"

	"github.com/gin-gonic/gin"
)

var Router *gin.Engine

func init() {
	Router = gin.Default()
	Router.Use(cors())
	v1 := Router.Group("v1")
	{

		// email
		emailApi := new(apiV1.EmailApi)
		v1.POST("/email", emailApi.Send)

		// user
		user := v1.Group("/user")
		userApi := new(apiV1.UserApi)
		user.POST("", userApi.Register)
		user.POST("/login", userApi.Login)
		user.PUT("/login", userApi.ResetPassword)
		user.Use(JWT())
		user.PUT("", userApi.UpdatePassword)
		user.GET("", userApi.UserInfo)

		// spaces
		v1.GET("/spaces")

		// providers
		v1.GET("/providers")
	}

}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "*")
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

var jwtService service.JWTService

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.Request.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			c.Abort()
			c.JSON(http.StatusOK, result.Result{
				Code: result.UserTokenInvalid,
				Msg:  "token invalid, please login again",
			})
			return
		}
		token := strings.TrimPrefix("Bearer ", auth)
		uid, newToken, err := jwtService.Validate(token)
		if err != nil {
			c.Abort()
			c.JSON(http.StatusOK, result.Result{
				Code: result.UserTokenExpired,
				Msg:  "token expired, please login again",
			})
			return
		}
		if newToken != "" {
			c.Header("new-token", newToken)
			c.Request.Header.Set("Authorization", "Bearer "+newToken)
		}
		c.Set("uid", uid)
		c.Next()
	}
}
