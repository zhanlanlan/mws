package controllers

import (
	"context"
	"fmt"
	"mws/model"

	"github.com/expectedsh/go-sonic/sonic"
	"github.com/gin-gonic/gin"
)

// QueryRecipe 根据菜名查找菜谱
func QueryRecipe(c *gin.Context) {

	cai, ok := c.GetQuery("cai")
	if !ok {
		c.JSONP(400, gin.H{
			"code":    "400",
			"message": "no query found",
		})
		return
	}

	search, err := sonic.NewSearch("127.0.0.1", 27016, "SecretPassword")
	if err != nil {
		c.JSONP(500, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	results, err := search.Query("mws", "recipe_list", cai, 10, 0)
	if err != nil {
		c.JSONP(500, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	}

	if err := model.O.Ping(context.TODO(), nil); err != nil {
		c.JSONP(500, gin.H{
			"code":    500,
			"message": "internal server error",
		})
		return
	} else {
		fmt.Println("success connect to mongodb ")
	}

	c.JSON(200, gin.H{
		"code":    200,
		"message": "success",
		"data":    results,
	})
}
