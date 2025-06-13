FROM golang:1.23.1-alpine as builder
ENV GOPROXY=https://goproxy.cn
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o FindCode cmd/main.go


FROM alpine:latest
RUN mkdir -p /app
WORKDIR /app

# 安装 git（使用 apk）
RUN apk update && apk add --no-cache git python3

COPY --from=builder /build/etc etc/
COPY --from=builder /build/rule rule/
COPY --from=builder /build/run.sh .
COPY --from=builder /build/FindCode .

ENTRYPOINT ["/bin/sh", "/app/run.sh"]