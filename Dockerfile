# 使用官方的 Golang 运行时作为父镜像
FROM golang:1.21-alpine AS build-env

# 设置工作目录
WORKDIR /app

# 将当前目录内容复制到容器的 /app 下
COPY . ./

# 构建
RUN go build -o hello-world-temporal

# 使用一个轻量级的镜像作为基础镜像
FROM alpine:latest

# 从 build-env 阶段复制二进制文件到 /usr/local/bin/
COPY --from=build-env /app/hello-world-temporal /usr/local/bin/hello-world-temporal

# 设置容器启动时执行的命令
#CMD ["hello-world-temporal"]