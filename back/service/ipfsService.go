package service

import (
	"errors"
	"io/ioutil"
	"strings"

	shell "github.com/ipfs/go-ipfs-api"
)

/**
 * @Description: 从IPFS获取文件
 * @param cid IPFS文件CID
 * @param ipfsServiceUrl IPFS服务地址
 * @return string 文件内容
 * @return error 错误信息
 */
func HandleGetIPFSFile(cid, ipfsServiceUrl string) (string, error) {

	var IPFSNodeAddr = "http://47.113.204.64:5001"
	if ipfsServiceUrl != "" {
		IPFSNodeAddr = ipfsServiceUrl
	}

	// 创建一个sharness节点
	sh := shell.NewShell(IPFSNodeAddr)
	if sh == nil {

		return "", errors.New("IPFS shell连接失败")
	}
	// 从URL参数中获取CID
	// cid := "QmRGpFDBtctKnFTDwLTVY8yocgvCE1M8J2e2E3g5DEvvTt"

	// 从IPFS获取文件
	rc, err := sh.Cat(cid)
	if err != nil {

		return "", errors.New("从IPFS获取文件失败:" + err.Error())
	}
	body, err := ioutil.ReadAll(rc)

	if err != nil {

		return "", errors.New("读取文件失败:" + err.Error())
	}

	// 关闭阅读
	defer rc.Close()
	return string(body), nil

}

func HandleUploadIPFSFile(fileContent, ipfsServiceUrl string) (string, error) {
	var IPFSNodeAddr = "http://47.113.204.64:5001"
	if ipfsServiceUrl != "" {
		IPFSNodeAddr = ipfsServiceUrl
	}
	sh := shell.NewShell(IPFSNodeAddr)
	if sh == nil {

		return "", errors.New("IPFS shell连接失败")
	}
	// 将文件内容上传到IPFS
	cid, err := sh.Add(strings.NewReader(fileContent))
	if err != nil {

		return "", errors.New("上传文件到IPFS失败:" + err.Error())
	}
	return cid, nil

}
