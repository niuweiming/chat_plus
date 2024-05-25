# 使用官方的Go语言运行时作为基础镜像
FROM golang:latest AS builder

# 维护者信息
LABEL authors="2695062156@qq.com"

# 设置工作目录
WORKDIR /app

# 将宿主机的当前目录下的所有文件复制到容器的/app目录下
COPY . /app

# 设置Go代理和增加超时时间
ENV GOPROXY=https://goproxy.io,direct
ENV GONOSUMDB=*

# 构建Go应用
RUN go mod download && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# 使用另一个基础镜像，用于运行环境，这里依然使用Ubuntu但不含Go构建环境，减小最终镜像体积
FROM ubuntu:latest

# 复制二进制文件到新镜像
COPY --from=builder /app/main /app/

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV NODE_ENV=production

# 暴露应用需要的端口
EXPOSE 5768

# 容器启动命令
CMD ["./main"]
