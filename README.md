# ChainQA4chainmaker

面向腾讯云-长安链（V2.2版本）的ChainQA（区块链审计与查询平台）



## 项目架构说明

- back：后端代码，主要处理IPFS加解密、转换Excel数据、封装请求等基础操作
- front：前端代码，页面展示
- contract：部署在长安链（腾讯云 https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1）的合约代码，**请注意此合约代码仅适配长安2.1~2.2版本**。
- doc：一些可用于测试的文档（excel）
- tencent-chainmaker：用于链通后端与腾讯云区块链的SDK。**需进行配置**。
- docker-compose.yml：用于启动docker的docker编排文件。

### 后端（Back）

👋如只是**<u>单纯体验后端</u>**，请按如下方式运行：

安装依赖包：

```bash
go mod vendor
```

运行：

```bash
go run main.go
```

后端监听9000端口。

👋如想部署**<u>整个</u>项目**，请参见docker章节。

### 前端（front）

👋如只是**<u>单纯体验前端</u>**，请按如下方式运行：

需提前安装node、yarn。

运行(前端调试环境下)：

```bash
yarn dev
```

👋如想部署**<u>整个</u>项目**，请参见docker章节。

### 区块链智能合约

位于contract文件夹下。

其中`tencentChainqaContractV221demo01.7z`为编译后的版本。`src`为源代码文件夹。

> 请注意本文件的合约代码仅适配长安2.1~2.2版本，可部署至[腾讯云长安链](https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1)

合约名建议为：`tencentChainqaContractV221demo01`（否则需要在前端更改合约名）

![合约方法](./doc/img/contract-method.png)

### tencent-chainmaker

主要是后端与 [腾讯云长安链](https://console.cloud.tencent.com/tbaas/chainmaker/chain/chainmaker-demo/basicInfo?chainId=1&demo=1) <u>**交互（调用合约）**</u>的系统。是对腾讯云提供的区块链交互的SDK的二次封装

本内容需要进行配置：

1. 前往[腾讯云API](https://console.cloud.tencent.com/cam/capi) 页面，申请密钥，保存SecretId和SecretKey
2. 进入`/tencent-chainmaker/conf/config.ini`文件，填入相应部分

> 密钥是高危部分！不能以任何形式公开！需妥善保管

> 为什么需要这个part，因为整体的链路是这样的
>
> 前端———后端——–tencent-chainmaker———长安链

### docker

下面讲解如何配置docker，打包为docker运行。

【第一步：打包前端】（可选：若对前端源代码进行了修改，需要按照本步骤重新打包）

- 安装[node](https://blog.csdn.net/muzidigbig/article/details/80493880)、[yarn](https://blog.csdn.net/qq_63055262/article/details/135776853)两个前端工具
- 切换到前端目录（`/front`），运行`yarn build`
- 稍等一会，就会在`/front/dist`文件夹下生成前端打包文件（文件中原`/front/dist`目录下已经存放了前端打包文件，可直接使用。若对源代码进行了修改，需要按照本步骤重新打包）

【第二步：前端打包文件放到后端文件夹】（可选：若对前端源代码进行了修改，需要按照本步骤执行，否则可以跳过，因为目前后端有这样一份代码）

- 将`/front/dist`整体复制到`/back/dist`，其他什么都不用干

【第三步：配置tencent-chainmaker】

- 请按本文档tencent-chainmaker部分进行配置SecretId和SecretKey

【第四步：安装区块链智能合约】

- 请按本文档区块链智能合约部分配置

【第五步：容器启动！】

切换到根目录下，运行：

```bash
docker-compose up -d
```

这样就会自动启动根目录下的`docker-compose.yml`，`yml`中又定义了两个容器（一个是后端，一个是tencent-chainmaker。第二步中前端已经打包放入后端了，所以没有前端容器）。然后会分别按照`back/Dockerfile`和`tencent-chainmaker/Dockerfile`进行镜像装载，然后进行启动。

> 如果没有更改默认配置，后端将运行在9000，tencent-chainmaker运行在9001端口

【第六步：运行】

打开浏览器，输入`localhost:9000`即可看到页面。

> 注：若需要将后端端口号进行修改，请修改`/back/conf/conf.ini`中的port

【第七步：修改】

![](./doc/img/modify-in-front.png)

你需要修改后端容器的API地址，请修改第一条。注意只需要修改`主机IP:端口号`的部分（即47.113.204.64:9000），主机IP为服务器IP（若为本地电脑请使用localhost），端口号如果你没有在配置文件`/back/conf/conf.ini`做任何修改，请沿用9000。

你可能需要修改IPFS容器地址。本文档没有提到IPFS容器的安装配置，如果需要，可以参考[我的这篇文档](https://www.yuque.com/jjq0425/pku/cm112pwu470v3q9n)。安装IPFS容器后，设置IP地址即可。

你需要修改tencent-chainmaker容器的API地址，请修改第三条。注意只需要修改`主机IP:端口号`的部分（即47.113.204.64:9000），主机IP为服务器IP（若为本地电脑请把47.113.204.64改为`host.docker.internal`或`172.17.0.1`或`localhost`，由于不同OS下的docker的本地环路地址不同，请把上面三者遍历尝试一遍）

最后一个为合约名称，是你在腾讯云区块链上填写的合约的名称。

> 注意：修改后点击保存后仅当次浏览有效！**若刷新浏览器上述配置需要重新设置**！
>
> 可以配置好后，修改前端`\front\src\stores\api.js`然后重新打包。

