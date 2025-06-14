package main

import "ScanIDOR/internal/pkg/server"

// 使用gin 搭建后段
func main() {
	server.NewServer().Run()
}
