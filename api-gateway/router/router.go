package router

import (
	"api-gateway/internal/handler"
	"api-gateway/middleware"
	"github.com/gin-gonic/gin"
	"net/http"
)

func NewRouter(services map[string]interface{}) *gin.Engine {
	engine := gin.Default()
	engine.Use(middleware.NoCache(),
		middleware.Options(),
		middleware.Secure(),
		middleware.InitMiddleware(services))

	engine.GET("ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, "ok")
	})

	v1 := engine.Group("/api/v1")
	{
		// user模块路由
		v1.POST("/user/register", handler.UserRegister)
		v1.POST("/user/login", handler.UserLogin)

		// task模块路由
		v1.POST("/task", handler.CreateTask)
	}
	return engine
}
