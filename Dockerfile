# 第一阶段：构建应用
FROM golang:1.23.1-alpine as builder

# 设置环境变量
ENV GOPROXY=https://goproxy.cn
ENV CGO_ENABLED=1
ENV GOOS=linux
ENV GOARCH=amd64

# 安装必要的依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

WORKDIR /build
COPY . .

# 下载依赖并构建应用
RUN go mod tidy
RUN go build -o FindCodeServer cmd/server.go

# 第二阶段：创建运行时镜像
FROM alpine:latest

# 安装 SQLite 运行时依赖
RUN apk add --no-cache sqlite-libs

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /build/etc etc/
COPY --from=builder /build/rule rule/
COPY --from=builder /build/web web/
COPY --from=builder /build/FindCodeServer .

ENTRYPOINT ["./FindCodeServer"]