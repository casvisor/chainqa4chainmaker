package controller

import (
	"chainqa_offchain_demo/models"
	"chainqa_offchain_demo/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

func DownloadIPFSFileHandler(c *gin.Context) {
	type UploadFileDTO struct {
		ApiUrl ApiUrlDTO `json:"apiUrl"` // API地址
		Cid    string    `json:"cid"`    // IPFS文件CID
	}

	var uploadFileDTO UploadFileDTO
	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&uploadFileDTO); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err.Error())
		return
	}

	// 调用service层的HandleGetIPFSFile方法，传入apiUrl和cid参数
	fileContent, err := service.HandleGetIPFSFile(uploadFileDTO.Cid, uploadFileDTO.ApiUrl.IpfsServiceUrl)
	if err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "获取文件失败", err.Error())
		return
	}

	models.ResponseOK(c, "获取IPFS文件（密文）成功", fileContent)
}

func TryDecryptFileHandler(c *gin.Context) {
	type UploadFileAndTryDecrptDTO struct {
		ApiUrl ApiUrlDTO `json:"apiUrl"` // API地址
		Cid    string    `json:"cid"`    // IPFS文件CID
		AesKey string    `json:"aesKey"` // AES密钥
	}

	var uploadFileAndTryDecrptDTO UploadFileAndTryDecrptDTO
	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&uploadFileAndTryDecrptDTO); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}

	// 调用service层的HandleGetIPFSFile方法，传入apiUrl和cid参数
	fileCiperContent, err := service.HandleGetIPFSFile(uploadFileAndTryDecrptDTO.Cid, uploadFileAndTryDecrptDTO.ApiUrl.IpfsServiceUrl)
	if err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "获取文件失败", err.Error())
		return
	}

	fileContent, err := service.AesDecrypt(fileCiperContent, uploadFileAndTryDecrptDTO.AesKey)
	if err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "解密文件失败", err.Error())
		return
	}
	models.ResponseOK(c, "解密文件（明文）成功", fileContent)

}
