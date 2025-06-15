package main

import (
	"ScanIDOR/internal/pkg/env"
	"ScanIDOR/internal/pkg/server/app"
	"ScanIDOR/internal/pkg/server/config"
	"ScanIDOR/utils/util"
	"fmt"
	"os"
)

// 使用gin 搭建后段
func main() {
	conf := config.Config{
		DbConfig: &config.SqliteConfig{
			FilePath: "data.db",
		},
		GinConfig: &config.GinConfig{
			ListenOn: ":8080",
			Mode:     "debug",
		},
	}

	a, err := app.NewApp(&conf)
	if err != nil {
		panic(err)
	}
	initToken()
	a.Run()
}

func initToken() {
	token := util.GenerateToken()
	env.Env["token"] = token
	fileName := "token.txt"
	file, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("create file %s failed.\n", fileName)
	}
	defer file.Close()
	_, err = file.WriteString(token)
}
