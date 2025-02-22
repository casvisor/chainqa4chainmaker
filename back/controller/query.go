package controller

import (
	"chainqa_offchain_demo/models"
	"chainqa_offchain_demo/service"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func QueryDataHandler(c *gin.Context) {
	type QueryDataDTO struct {
		Uid       string    `json:"uId"`       // 用户ID
		QueryItem string    `json:"queryItem"` // 查询项
		ApiUrl    ApiUrlDTO `json:"apiUrl"`    // API地址
		// FilePoses []string  `json:"filePoses"` // 文件位置（废弃，直接从QueryItem解析）
	}

	var queryDataDTO QueryDataDTO
	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&queryDataDTO); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}
	queryDataDTO.Uid = strings.TrimSpace(queryDataDTO.Uid)             // 去除空格
	queryDataDTO.QueryItem = strings.TrimSpace(queryDataDTO.QueryItem) // 去除空格

	// 提取filePos
	// 解析 QueryItem 字段中的 JSON
	type QueryItemDTO struct {
		FilePos [][]string `json:"filePos"` //  filePos 是一个二维数组
	}
	var queryItemDTO QueryItemDTO
	if err := json.Unmarshal([]byte(queryDataDTO.QueryItem), &queryItemDTO); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "QueryItem 格式错误", err)
		return
	}
	// 提取 filePos 字段
	filePos2Dim := queryItemDTO.FilePos
	if len(filePos2Dim) == 0 {
		models.ResponseError400(c, http.StatusBadRequest, "filePos 不能为空", nil)
		return
	}
	// 将filePos变为一维数组
	FilePoses := make([]string, 0) // 初始化切片
	for _, pos := range filePos2Dim {
		for _, p := range pos {
			FilePoses = append(FilePoses, p)
		}
	}

	// 查询每个文件位置的数据，用map[string]string存储
	fileDataMap := make(map[string]string)
	for _, filePos := range FilePoses {
		fileCiperData, err := service.HandleGetIPFSFile(filePos, queryDataDTO.ApiUrl.IpfsServiceUrl)
		if err != nil {
			models.ResponseError400(c, http.StatusBadRequest, "请求IPFS错误:"+err.Error(), err)
			return
		}
		// 等待2s
		time.Sleep(2000 * time.Millisecond) // 延时2s
		// 从区块链中获取AES密钥
		aesKey, err := service.GetAesKeyFromBlockchain(queryDataDTO.ApiUrl.ContractName, queryDataDTO.ApiUrl.ChainServiceUrl, filePos)
		if err != nil {
			models.ResponseError400(c, http.StatusBadRequest, "获取AES密钥失败:"+err.Error(), err)
			return
		}
		// 解密
		fileData, err := service.AesDecrypt(fileCiperData, aesKey) // 解密文件
		if err != nil {
			models.ResponseError400(c, http.StatusBadRequest, "解密失败:"+err.Error(), err)
			return
		}
		fileDataMap[filePos] = fileData
	}

	// 调用查询函数
	queryResult, code := service.GetQueryResult(queryDataDTO.QueryItem, fileDataMap) // 返回查询结果

	time.Sleep(1000 * time.Millisecond)                                                                                                                               // 延时1s
	err := service.UpdateQueryLog(queryDataDTO.ApiUrl.ContractName, queryDataDTO.ApiUrl.ChainServiceUrl, queryDataDTO.Uid, queryDataDTO.QueryItem, code, queryResult) // 上链查询日志
	if err != nil {
		fmt.Println("上链查询日志失败", err)
	}
	if code == -1 {
		models.ResponseOK(c, "查询完成，但出现错误", queryResult)
		return
	} else {
		models.ResponseOK(c, "查询成功", queryResult)
	}

	// TODO: 缓存

}
