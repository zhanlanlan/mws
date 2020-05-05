package router

import (
	"mws/controllers"

	"github.com/gin-gonic/gin"
)

// Route 路由定义函数
func Route() *gin.Engine {
	router := gin.Default()

	router.GET("/recipe/query", controllers.QueryRecipe)

	return router
}
