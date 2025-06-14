FROM golang:1.23.1-alpine as builder
ENV GOPROXY=https://goproxy.cn
WORKDIR /build
COPY . .
RUN go mod tidy
RUN go build -o FindCodeServer cmd/server.go


FROM alpine:latest
RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /build/etc etc/
COPY --from=builder /build/rule rule/
COPY --from=builder /build/templates templates/
COPY --from=builder /build/FindCodeServer .

ENTRYPOINT ["./FindCodeServer"]