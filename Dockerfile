# 构建阶段：使用固定版本的 Go 镜像
FROM golang:1.24-alpine AS builder

# 设置环境变量
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct \
    CGO_ENABLED=0 \
    GOOS=linux

# 设置工作目录（使用绝对路径避免混淆）
WORKDIR /app

# 复制依赖文件（利用 Docker 缓存层）
COPY go.mod go.sum ./
# 带重试机制的依赖下载
RUN for i in $(seq 1 10); do \
        go mod download && break || \
        { echo "requirement download failed， retry... $i ..."; sleep 10; }; \
    done

# 复制源代码并构建
COPY . .
RUN go build -ldflags="-s -w" -o bin/agent . \
    && chmod +x bin/lanaya

FROM alpine:3.19
WORKDIR /app
COPY --from=builder  /app/bin/agent .
USER root
CMD ["./agent"]