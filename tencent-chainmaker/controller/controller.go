package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tencent-chainmaker/models"
	"tencent-chainmaker/setting"

	"github.com/gin-gonic/gin"

	// 腾讯云区块链服务
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"

	tbaas "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tbaas/v20180416"
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

type chainDTO struct {
	ContractName string                 `json:"contractName"`
	MethodName   string                 `json:"methodName"`
	Args         map[string]interface{} `json:"args"`
}

// /exec
func ExecChain(c *gin.Context) {
	fmt.Println("exec chain")
	var req chainDTO

	// 绑定JSON数据到结构体
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ResponseError400(c, http.StatusBadRequest, "请求格式错误", err)
		return
	}

	credential := common.NewCredential(
		setting.Conf.GetSecretId(),
		setting.Conf.GetSecretKey(),
	)

	// 将map转换为JSON字符串
	jsonData, err := json.Marshal(req.Args)
	if err != nil {
		models.ResponseError500(c, http.StatusInternalServerError, "JSON转换失败", err)
		return
	}

	fmt.Printf("jsonData: %v\n", string(jsonData))

	// 实例化一个client选项，可选的，没有特殊需求可以跳过
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tbaas.tencentcloudapi.com"
	// 实例化要请求产品的client对象,clientProfile是可选的
	client, _ := tbaas.NewClient(credential, "ap-beijing", cpf)

	// 实例化一个请求对象,每个接口都会对应一个request对象
	request := tbaas.NewInvokeChainMakerDemoContractRequest()

	// TODO: 设置请求参数，可能需要修改。目前匹配腾讯云测试网络
	request.ClusterId = common.StringPtr("chainmaker-demo")   // 网络ID（写死）
	request.ChainId = common.StringPtr("chain_demo")          // 链ID(写死)
	request.ContractName = common.StringPtr(req.ContractName) // 合约名
	request.FuncName = common.StringPtr(req.MethodName)       // 合约方法名
	request.FuncParam = common.StringPtr(string(jsonData))    // 合约方法参数
	request.AsyncFlag = common.Int64Ptr(0)                    // 同步响应标志，1:异步响应 0:同步响应

	// 返回的resp是一个InvokeChainMakerDemoContractResponse的实例，与请求对象对应

	response, err := client.InvokeChainMakerDemoContract(request)
	fmt.Printf("response: %v\n", response)
	fmt.Printf("err: %v\n", err)
	if _, ok := err.(*errors.TencentCloudSDKError); ok {
		models.ResponseError500(c, http.StatusInternalServerError, "腾讯云服务异常", err)
		return
	}
	if err != nil {
		models.ResponseError500(c, http.StatusInternalServerError, "腾讯云服务异常2", response)
		return

	}
	// 输出json格式的字符串回包
	models.ResponseOK(c, "success", response)
}
