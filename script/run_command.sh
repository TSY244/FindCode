#!/bin/bash
cd .. # 因为需要etc/rule 等默认配置文件
go build -ldflags="-s -w" -o ./FindCodeCommand ./cmd/command/main.go  && ./FindCodeCommand  $@