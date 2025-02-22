package controller

import (
	"chainqa_offchain_demo/models"
	"chainqa_offchain_demo/service"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ApiUrlDTO struct {
	IpfsServiceUrl  string `json:"ipfsServiceUrl"`  // IPFS服务地址
	ContractName    string `json:"contractName"`    // 合约名
	ChainServiceUrl string `json:"chainServiceUrl"` // 链服务地址
}

// hello
func HelloHandler(c *gin.Context) {
	// 延时300ms
	// time.Sleep(30000 * time.Millisecond)
	models.ResponseOK(c, "success", "hello world")
}

// upload

func GetAesKeyHandler(c *gin.Context) {
	type GetAesKeyDTO struct {
		KeyNum int `json:"keyNum"`
	}
	var req GetAesKeyDTO

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}

	// 检查keyNum是否为正数，如果不是，可以设置默认值或者返回错误
	if req.KeyNum <= 0 {
		req.KeyNum = 1 // 设置默认值
	}

	// 调用service层的GetAesKey方法，传入keyNum参数
	KeyArr := make([]string, 0)
	for i := 0; i < req.KeyNum; i++ { // 使用var i声明循环变量
		key, err := service.GenerateAESKey()
		if err != nil {
			models.ResponseError400(c, http.StatusBadRequest, "获取AES密钥失败", err)
			return
		}
		KeyArr = append(KeyArr, key)
	}

	// 假设你需要返回KeyArr，以下是返回的示例代码
	models.ResponseOK(c, fmt.Sprintf("获取%d个AES密钥成功", req.KeyNum), KeyArr)
}

func UploadFileHandler(c *gin.Context) {
	type UploadFileDTO struct {
		ApiUrl      ApiUrlDTO `json:"apiUrl"`      // API地址
		AesKey      string    `json:"aesKey"`      // AES密钥
		FileContent string    `json:"fileContent"` // 文件内容
		FileName    string    `json:"fileName"`
		Uid         string    `json:"uId"`
	}
	var req UploadFileDTO

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}
	req.Uid = strings.TrimSpace(req.Uid)                 // 去除空格
	req.FileName = strings.TrimSpace(req.FileName)       // 去除空格
	req.FileContent = strings.TrimSpace(req.FileContent) // 去除空格
	req.AesKey = strings.TrimSpace(req.AesKey)           // 去除空格

	// 1. 使用AES密钥加密文件内容
	cipherText, err := service.AesEncrypt(req.FileContent, req.AesKey)

	type UploadFileErrorVO struct {
		FileName string `json:"fileName"` // 文件名
		Error    string `json:"error"`    // 错误信息
	}
	var errorVO UploadFileErrorVO
	errorVO.FileName = req.FileName
	if err != nil {
		errorVO.Error = err.Error()
		models.ResponseError400(c, http.StatusBadRequest, fmt.Sprintf("加密文件内容失败: %s", req.FileName), errorVO)
		return
	}

	// 2. 将加密后的文件内容上传到IPFS
	cid, err := service.HandleUploadIPFSFile(cipherText, req.ApiUrl.IpfsServiceUrl) // 调用HandleUploadIPFSFile方法上传文件
	if err != nil {
		errorVO.Error = err.Error()
		models.ResponseError400(c, http.StatusBadRequest, fmt.Sprintf("上传文件到IPFS失败: %s", req.FileName), errorVO)
		return
	}

	// 3. 上传数字信封
	// 从区块链中获取公钥
	publicKey, err := service.GetPublicKeyFromBlockchain(req.ApiUrl.ContractName, req.ApiUrl.ChainServiceUrl)
	if err != nil {
		errorVO.Error = err.Error()
		models.ResponseError400(c, http.StatusBadRequest, "获取公钥失败", errorVO)
		return
	}
	Envelope, err := service.RSAEncryptAndReturnEnvelop(req.ApiUrl.ContractName, req.AesKey, publicKey)
	if err != nil {
		errorVO.Error = err.Error()
		models.ResponseError400(c, http.StatusBadRequest, fmt.Sprintf("生成数字信封失败: %s", req.FileName), errorVO)
		return
	}
	// 等待1s
	time.Sleep(1000 * time.Millisecond) // 延时1s
	// 将数字信封上传到区块链
	err = service.UploadEnvelopeToBlockchain(req.ApiUrl.ContractName, req.ApiUrl.ChainServiceUrl, Envelope, cid, req.Uid)
	if err != nil {
		errorVO.Error = err.Error()
		models.ResponseError400(c, http.StatusBadRequest, fmt.Sprintf("上传数字信封到区块链失败: %s", req.FileName), errorVO)
		return
	}

	// 4. 返回文件上传结果
	type UploadFileSuccessVO struct {
		FileName string `json:"fileName"` // 文件名
		Pos      string `json:"pos"`      // ipfs地址
	}

	var successVO UploadFileSuccessVO
	successVO.FileName = req.FileName
	successVO.Pos = cid
	models.ResponseOK(c, fmt.Sprintf("上传文件%s到IPFS成功", req.FileName), successVO)

}
