# Tencent-chainmaker 模块介绍

# 第一部分：腾讯云长安链（V2.2.1）交互模块部署

## 说明

### 模块作用

本模块为用户后端与腾讯云长安链的**交互中间件**。是对腾讯云长安链SDK的二次封装，用于对[腾讯云长安链](https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1)进行合约调用、上链、查询

### 使用方法

1. 前往腾讯云申请API ID和密钥（secretID、secretKey）。网址：[API密钥管理](https://console.cloud.tencent.com/cam/capi)

   > 注意：密钥为高危内容，不能公开！
   >
2. 修改 `/conf/config_template.ini`的腾讯云API密钥

   > 如果修改配置参数的port，那么下面第二部分对应的端口号也需要跟着改变
   >
3. 复制 `/conf/config_template.ini`文件，粘贴到**同级**目录（`/conf`）下，重命名为 `config.ini`。后续程序将读取 `config.ini`文件作为配置参数。
4. 输入下面指令启动程序

```bash
cd 项目根目录
go run main.go
```

然后，应用就可以使用交互模块提供的接口与长安链交互了。

接口详见第二部分。

# 第二部分：腾讯云长安链（V2.2.1）交互接口

## 基础信息

Base URLs:

* `<a href="http://localhost:9000/api">`开发环境: http://localhost:9000/api`</a>`

### 接口1：POST 链通测试

- 请求方式：POST
- 请求网址后缀：/hello（此前需要加上Base URLs，请见第二部分开头的说明）

返回类似 `hello world`的内容，可供测试是否部署正确。本接口不与腾讯云区块链交互。

### 接口2：POST 进行交易

> 区块链的query和invoke均支持本接口

#### 请求体

- 请求方式：POST
- 请求网址后缀：/exec （此前需要加上Base URLs，请见第二部分开头的基础信息部分）
- 请求参数（Body）：

```json
{
  "contractName": "string",
  "methodName": "string",
  "args": {}
}
```

| 名称            | 位置 | 类型         | 必选 | 说明   |
| --------------- | ---- | ------------ | ---- | ------ |
| body            | body | object       | 否   | none   |
| » contractName | body | string       | 是   | 合约名 |
| » methodName   | body | string       | 是   | 方法名 |
| » args         | body | object¦null | 是   | 参数   |

args请自行填写所需参数，以JSON格式，**请求示例**如下：

```JSON
{
    "contractName": "ChainMakerDemo",
    "methodName": "save",
    "args": {
        "key":"test",
        "field":"test",
        "value":"test"
    }
}
```

#### 返回结果

> 200 Response

举例：

```json
{
    "code": 0,
    "data": {
        "Response": {
            "Result": {
                "Code": 0,
                "CodeMessage": "Success",
                "TxId": "ce549d273c7b474396d696e438b23fb7c4d6206e868e48da8ede035f9bff3db0",
                "GasUsed": 11451,
                "Message": "Success",
                "Result": "U3VjY2Vzcw=="
            },
            "RequestId": "4d8bcf17-4a00-420b-8b85-e3aa5c2d8282"
        }
    },
    "msg": "success"
}
```

| 状态码 | 状态码含义                                           | 说明 | 数据模型 |
| ------ | ---------------------------------------------------- | ---- | -------- |
| 200    | [OK](https://tools.ietf.org/html/rfc7231#section-6.3.1) | none | Inline   |

其余结果请参考返回示例。其中Result即我们上链返回的数据的**base64编码结果**

> 腾讯长安链不区分 上链（invoke） OR 查询（query）。因此统一整合到这一个接口

## 使用方法

比如你的应用想要调用 `ChainMakerDemo`合约的save方法，传参key，field，value。你只需要向：`http://localhost:9000/api/exec` 发送一个POST请求（请求体同请求示例）。那么你就可以收到类似于返回结果示例的内容。
