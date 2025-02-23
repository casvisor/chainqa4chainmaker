<h1 id="rvDSA">IPFS安装下载</h1>

+ <font style="color:rgb(38, 38, 38);">在你的服务器中，创建一个空白文件夹，为IPFS根目录
+ 在根目录下创建一个`ipfsDataAndExport`文件夹与一个`bash`文件夹，在`bash`文件夹下编写如下 shell 脚本，**请注意下面要修改的部分**。</font>

```shell
#!/bin/bash

# ====以下内容可能需要修改=========
# 设置跨域的命令
CORS_ORIGIN='["http://*****:5001", "http://localhost:3000", "http://127.0.0.1:5001", "https://webui.ipfs.io"]'
# 【tip1】若为服务器，请修改第一个IP为服务器公网IP。否则删掉第一个！
CORS_METHODS='["PUT", "POST"]'


# 定义目录变量
base_dir="../ipfsDataAndExport"
# 【tip2】请以本脚本的位置为定位，按相对位置创建上面ipfsDataAndExport文件夹，下面文件夹不需要创建，会自动化创建


# 容器名
CONTAINER_NAME=ipfs_host
# 【tip3】容器名，请按照实际需求修改。

# ====以上内容可能需要修改=========

ipfs_data_dir="$base_dir/ipfs_data"
ipfs_export_dir="$base_dir/ipfs_export"


# 检查并删除现有的文件夹
if [ -d "$ipfs_data_dir" ]; then
  echo "删除旧ipfs_data"
  rm -rf "$ipfs_data_dir"
fi

if [ -d "$ipfs_export_dir" ]; then
  echo "删除旧ipfs_export"
  rm -rf "$ipfs_export_dir"
fi

# 等待容器启动
echo "等待删除程序，等待3s..."
sleep 3

# if [ $# -eq 0 ]; then
#   echo "未传入IPFS Docker连接网络名称参数。本脚本要求传入IPFS Docker连接网络名称参数"
#   echo "提示：因为需要保证区块链背书节点也能访问到ipfs，所以需要ipfs和区块链使用同一个网络，所以需要设置网络名称。建议先启动区块链。确认容器的网络配置是否正确!"
#   exit 1
# fi


#---------------------------------------------------------------





# 创建新的文件夹
echo "(1) 创建新ipfs_export和ipfs_data"
mkdir -p "$ipfs_data_dir"
mkdir -p "$ipfs_export_dir"

# 输出结果
if [ -d "$ipfs_data_dir" ] && [ -d "$ipfs_export_dir" ]; then
  echo -e "\033[0;32m 创建新ipfs_export和ipfs_data成功 \033[0m"
else
  echo "Failed to create folders."
  exit 1
fi



echo "(2) 启动ipfs_host容器..."
# export ipfs_host_NETWORK_NAME=$1
export ipfs_host_CONTAINER_NAME=$CONTAINER_NAME
docker-compose up -d

# 等待容器启动
echo "等待容器自行启动，等待10s..."
sleep 10

echo ">> 可以使用docker logs -f ipfs_host 自行查看容器是否启动成功"
echo ">> 若末尾出现Daemon is ready 说明启动成功"

# # 定义要查找的字符串
# STRING_TO_FIND="Daemon is ready"
# # 使用 docker logs -f 和 awk 来等待特定的日志消息
# docker logs -f ipfs_host | awk '/'"$STRING_TO_FIND"'/{print; exit}'

# # 检查 awk 的退出状态
# if [ $? -eq 0 ]; then
#   echo -e "\033[0;32m 容器启动成功 \033[0m"
#   # 在这里添加后续命令
# else
#   echo "错误：没有找到消息 '$STRING_TO_FIND' ，容器启动失败。"
#   exit 1
# fi


echo "(3) 开始设置跨域配置..."

# 进入容器内部并设置跨域
docker exec -it $CONTAINER_NAME sh -c "
    ipfs config --json API.HTTPHeaders.Access-Control-Allow-Origin '$CORS_ORIGIN' &&
    ipfs config --json API.HTTPHeaders.Access-Control-Allow-Methods '$CORS_METHODS'
  "

# 检查设置跨域命令是否成功
  if [ $? -eq 0 ]; then
  
    echo -e "\033[0;32m 跨域配置设置成功， \033[0m"
    echo "(4) 重启容器..."
    # 重启容器
    docker restart $CONTAINER_NAME
    echo -e "\033[0;32m 容器配置完成，可以访问 公网IP:5001/webui 或者 localhost:5001/webui 查看webUI界面 \033[0m"
  else
    echo "跨域配置设置失败。"
    exit 1
  fi
```

+ <font style="color:rgb(38, 38, 38);">同一目录下创建 docker-compose.yml</font>

```yaml
version: '2'

services:
  ipfs_host:
    container_name: "${ipfs_host_CONTAINER_NAME}"
    image: ipfs/kubo:latest
    volumes:
      - ../ipfsDataAndExport/ipfs_export:/export
      - ../ipfsDataAndExport/ipfs_data:/data/ipfs
    #command: /bin/bash   
    ports:
      - 4001:4001/udp
      - 8080:8080
      - 5001:5001
      - 8080:8080
```

+ <font style="color:rgb(38, 38, 38);">运行 shell 脚本</font>

```bash
./create-ipfs-docker.sh
```

<h3 id="RwMFw">Tips</h3>
1. **<font style="color:#7E45E8;">容器安装疑似出现问题？怎么判断容器是否安装成功？</font>**

<font style="color:rgb(38, 38, 38);">若在容器安装时出现问题（还没进行到</font>_<font style="color:rgb(38, 38, 38);">(3) 开始设置跨域配置...</font>_<font style="color:rgb(38, 38, 38);">）可以使用docker logs -f ipfs_host自行查看容器是否启动成功，若末尾出现Daemon is ready 说明启动成功（截图见下）</font>

```plain
Run 'ipfs id' to inspect announced and discovered multiaddrs of this node.
RPC API server listening on /ip4/0.0.0.0/tcp/5001
WebUI: http://127.0.0.1:5001/webui
Gateway server listening on /ip4/0.0.0.0/tcp/8080
Daemon is ready
```

<font style="color:rgb(38, 38, 38);">访问 服务器:5001/webui即可得到网页版</font>

:::danger
5001 端口切勿开启公网访问！可以限制仅某些 IP（如校内 IP）可以访问

:::



2. **<font style="color:#7E45E8;">相比版本 2，不在同一网络了，如何使得外部能够</font>****<u><font style="color:#7E45E8;">便捷</font></u>****<font style="color:#7E45E8;">访问呢？</font>**

Q：我在宿主机运行 IPFS docker，然后我也运行了另一个 docker B，如何在 docker B 中访问 IPFS

A：关键代码段

```go
const IPFSNodeAddr = "172.17.0.1:5001"
// 创建一个sharness节点
sh := shell.NewShell(IPFSNodeAddr)
```

即在不改变默认网络（docker B 没有进行天花乱坠的网路配置，目前正常使用区块链无问题）的条件下，172.17.0.1 即可访问宿主机的网路



