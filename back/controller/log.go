package controller

import (
	"chainqa_offchain_demo/models"
	"chainqa_offchain_demo/service"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func LogByUidHandler(c *gin.Context) {
	type LogByUid struct {
		Uid    string    `json:"uId"`    // 用户ID
		ApiUrl ApiUrlDTO `json:"apiUrl"` // API地址
	}

	var logByUid LogByUid
	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&logByUid); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}
	logByUid.Uid = strings.TrimSpace(logByUid.Uid) // 去除空格

	// 调用查询函数
	queryResult, err := service.GetAllQueryLogByUid(logByUid.ApiUrl.ContractName, logByUid.ApiUrl.ChainServiceUrl, logByUid.Uid) // 返回查询结果和错误码
	if err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "查询失败", err)
		return
	}

	models.ResponseOK(c, "查询成功", queryResult)
}

func LogByTimeRangeHandler(c *gin.Context) {
	type LogByTimestamp struct {
		ApiUrl    ApiUrlDTO `json:"apiUrl"`    // API地址
		StartTime string    `json:"startTime"` // 开始时间
		EndTime   string    `json:"endTime"`   // 结束时间
	}

	var logByTimestamp LogByTimestamp
	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&logByTimestamp); err != nil {

		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}
	// 调用查询函数
	queryResult, err := service.GetAllQueryLogByTimestamp(logByTimestamp.ApiUrl.ContractName, logByTimestamp.ApiUrl.ChainServiceUrl, logByTimestamp.StartTime, logByTimestamp.EndTime) // 返回查询结果和错误码
	if err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "查询失败", err)
		return
	}
	models.ResponseOK(c, "查询成功", queryResult)
}
