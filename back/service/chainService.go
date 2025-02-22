package service

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

type chainDTO struct {
	ContractName string                 `json:"contractName"`
	MethodName   string                 `json:"methodName"`
	Args         map[string]interface{} `json:"args"`
}

type ChainSuccessResponse struct {
	Code int `json:"code"`
	Data struct {
		Response struct {
			Result struct {
				Result string `json:"Result"`
			} `json:"Result"`
		} `json:"Response"`
	} `json:"data"`
	Msg string `json:"msg"`
}

// ExecBlockchain 执行区块链操作
func ExecBlockchain(chainServiceUrl string, jsonData []byte) (string, error) {
	// 创建HTTP客户端
	client := &http.Client{}

	// 创建请求对象
	req, err := http.NewRequest("POST", chainServiceUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", errors.New("创建请求失败" + err.Error())
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", err // 发送请求失败
	}
	defer resp.Body.Close() // 一定要关闭响应体

	// ======================= 处理响应 ======================
	// 定义结构体以匹配JSON响应

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.New("读取响应体失败" + err.Error())
	}
	if resp.StatusCode != http.StatusOK {
		return "", errors.New("请求失败，状态码: " + resp.Status)
	}

	// 解析响应体
	var chainSuccessResponse ChainSuccessResponse
	err = json.Unmarshal(body, &chainSuccessResponse)

	if err != nil {
		return "", errors.New("解析响应体失败" + err.Error())
	}

	decodedResult, err := base64.StdEncoding.DecodeString(chainSuccessResponse.Data.Response.Result.Result)
	if err != nil {
		return "", errors.New("解码失败" + err.Error())
	} else if chainSuccessResponse.Code != 0 {
		return "", errors.New("请求失败: " + string(decodedResult))
	}
	return string(decodedResult), nil

}

// ------------------------ 以下为业务函数 ------------------------
func GetPublicKeyFromBlockchain(contractName, chainServiceUrl string) (string, error) {
	// ====================== 构造响应 ======================
	if chainServiceUrl == "" {
		chainServiceUrl = "http://host.docker.internal:9001/tencent-chainapi/exec"
	}

	// 创建请求数据
	data := chainDTO{
		ContractName: contractName,
		MethodName:   "getPk",
		Args:         map[string]interface{}{
			// "file_prefix_path": contractName,
		},
	}

	// ======================= 发送请求 ======================
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("转换JSON失败" + err.Error())
	}

	return ExecBlockchain(chainServiceUrl, jsonData)

}

func UploadEnvelopeToBlockchain(contractName string, chainServiceUrl string, envelope string, cid string, uId string) error {
	// ====================== 构造响应 ======================
	if chainServiceUrl == "" {
		chainServiceUrl = "http://host.docker.internal:9001/tencent-chainapi/exec"
	}

	// 创建请求数据
	data := chainDTO{
		ContractName: contractName,
		MethodName:   "updateDataDigtalEnvelop",
		Args: map[string]interface{}{
			// "file_prefix_path": contractName,
			"uId":     uId,
			"pos":     cid,
			"envelop": envelope,
		},
	}

	// ======================= 发送请求 ======================
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("转换JSON失败" + err.Error())
	}

	_, err = ExecBlockchain(chainServiceUrl, jsonData)
	if err != nil {
		return err
	}
	return nil
}

func GetAesKeyFromBlockchain(contractName string, chainServiceUrl string, pos string) (string, error) {
	// ====================== 构造响应 ======================
	if chainServiceUrl == "" {
		chainServiceUrl = "http://host.docker.internal:9001/tencent-chainapi/exec"
	}

	// 创建请求数据
	data := chainDTO{
		ContractName: contractName,
		MethodName:   "getAesKey",
		Args: map[string]interface{}{
			"pos": pos,
		},
	}

	// ======================= 发送请求 ======================
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("转换JSON失败" + err.Error())
	}

	return ExecBlockchain(chainServiceUrl, jsonData)
}

// UpdateQueryLog 更新查询日志
func UpdateQueryLog(contractName string, chainServiceUrl string, uId string, queryItem string, queryStatus int, queryResult string) error {
	// ====================== 构造响应 ======================
	if chainServiceUrl == "" {
		chainServiceUrl = "http://host.docker.internal:9001/tencent-chainapi/exec"
	}

	queryStatusStr := strconv.Itoa(queryStatus)

	// 创建请求数据
	data := chainDTO{
		ContractName: contractName,
		MethodName:   "updateQueryLog",
		Args: map[string]interface{}{
			"uId":         uId,
			"queryItem":   queryItem,
			"queryStatus": queryStatusStr,
			"queryResult": queryResult,
		},
	}

	// ======================= 发送请求 ======================
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return errors.New("转换JSON失败" + err.Error())
	}

	_, err = ExecBlockchain(chainServiceUrl, jsonData)
	if err != nil {
		return err
	}
	return nil
}

func GetAllQueryLogByUid(contractName string, chainServiceUrl string, uId string) (string, error) {
	// ====================== 构造响应 ======================
	if chainServiceUrl == "" {
		chainServiceUrl = "http://host.docker.internal:9001/tencent-chainapi/exec"
	}

	// 创建请求数据
	data := chainDTO{
		ContractName: contractName,
		MethodName:   "getAllQueryLogByUid",
		Args: map[string]interface{}{
			"uId": uId,
		},
	}

	// ======================= 发送请求 ======================
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", errors.New("转换JSON失败" + err.Error())
	}

	return ExecBlockchain(chainServiceUrl, jsonData)
}

func GetAllQueryLogByTimestamp(contractName string, chainServiceUrl string, startTime string, endTime string) (string, error) {
	// ====================== 构造响应 ======================
	if chainServiceUrl == "" {
		chainServiceUrl = "http://host.docker.internal:9001/tencent-chainapi/exec"
	}

	// 创建请求数据
	data := chainDTO{
		ContractName: contractName,
		MethodName:   "getAllQueryLogByTimestamp",
		Args: map[string]interface{}{
			"startTime": startTime,
			"endTime":   endTime,
		},
	}

	// ======================= 发送请求 ======================
	// 将结构体转换为JSON
	jsonData, err := json.Marshal(data)
	if err != nil {

		return "", errors.New("转换JSON失败" + err.Error())
	}

	return ExecBlockchain(chainServiceUrl, jsonData)

}
