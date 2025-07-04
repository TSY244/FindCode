#!/bin/bash
cd .. # 因为需要etc/rule 等默认配置文件
go build -o FindCode cmd/main.go && ./FindCode -mode command $@