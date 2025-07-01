package fingerprint

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func GetProductPrint(productPath string) []string {
	// 遍历寻找 go.mod
	goModFilePath := getGoModPath(productPath)
	return getAllPrints(goModFilePath)
}

func getGoModPath(productPath string) string {
	absPath, err := filepath.Abs(productPath)
	if err != nil {
		return ""
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return ""
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "go.mod" {
			return filepath.Join(absPath, entry.Name())
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			subDirPath := filepath.Join(absPath, entry.Name())
			if found := getGoModPath(subDirPath); found != "" {
				return found
			}
		}
	}
	return "" // 未找到
}

func getAllPrints(modFilePath string) []string {
	file, err := os.Open(modFilePath)
	if err != nil {
		panic(err) // 因为是运行之前的检查所以这个地方很重要，因此可以直接抛出panic
	}
	defer file.Close()
	fileContent, err := io.ReadAll(file)

	prints := make([]string, 0)

	// 使用 WaitGroup 并发检测
	var wg sync.WaitGroup
	results := make(chan string)

	for printName, rules := range AllTasks {
		wg.Add(1)
		go func(pn string, rls []string) {
			defer wg.Done()
			for _, r := range rls {
				if strings.Contains(string(fileContent), r) {
					results <- pn
					return
				}
			}
		}(printName, rules)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	// 收集结果并去重
	found := make(map[string]bool)
	for p := range results {
		if !found[p] {
			prints = append(prints, p)
			found[p] = true
		}
	}

	return prints
}
