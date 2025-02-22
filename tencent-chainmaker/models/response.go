package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// 封装统一返回格式
func ResponseOK(c *gin.Context, msg string, data interface{}) {
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"msg":  msg,
		"data": data,
	})
}

// 400错误
func ResponseError400(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusBadRequest, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// 500错误
func ResponseError500(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusInternalServerError, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}

// 404错误
func ResponseError404(c *gin.Context, code int, msg string, data interface{}) {
	c.JSON(http.StatusNotFound, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
}
