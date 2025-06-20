package main

import (
	"ScanIDOR/internal/pkg/env"
	"ScanIDOR/internal/pkg/server/app"
	"ScanIDOR/internal/pkg/server/config"
	"ScanIDOR/utils/util"
	"fmt"
	"os"
	"strings"
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
	go clearDir()

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

func clearDir() {
	// 获取当前文件目录下的所有的文件夹
	files, err := os.ReadDir(".")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, file := range files {
		if file.IsDir() {
			// 删除文件夹
			if strings.HasPrefix(file.Name(), "git-clone-") {
				os.RemoveAll(file.Name())
			}
		}
	}

}
