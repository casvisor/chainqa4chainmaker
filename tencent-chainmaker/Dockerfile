FROM golang:alpine AS builder

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

    # 安装 ca-certificates
RUN apk add --no-cache ca-certificates

# 移动到工作目录：/build
WORKDIR /build

# 将代码复制到容器中
COPY . .

# 下载依赖信息
RUN go mod download

# 将我们的代码编译成二进制可执行文件 chainqa_offchain_demo
RUN go build -o chainqa_tencent_chainmaker .

# ###################
# # 接下来创建一个小镜像
# ###################
FROM scratch

# 从builder镜像中把配置文件拷贝到当前目录
COPY ./conf /conf

# 从builder镜像中把 ca-certificates 拷贝到当前目录
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# 从builder镜像中把/dist/app 拷贝到当前目录
COPY --from=builder /build/chainqa_tencent_chainmaker /
# 需要运行的命令
ENTRYPOINT ["/chainqa_tencent_chainmaker", "conf/config.ini"]