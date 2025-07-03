package fingerprint

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func getFilePath(productPath string, fileName map[string]struct{}) map[string]string {
	var result = make(map[string]string)
	absPath, err := filepath.Abs(productPath)
	if err != nil {
		return nil
	}

	entries, err := os.ReadDir(absPath)
	if err != nil {
		return nil
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if _, ok := fileName[entry.Name()]; ok {
				result[entry.Name()] = absPath + "/" + entry.Name()
			}
		}
	}

	for _, entry := range entries {
		if entry.IsDir() {
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}
			subDirPath := filepath.Join(absPath, entry.Name())
			if found := getFilePath(subDirPath, fileName); found != nil {
				result = mapAdd(result, found)
			}
		}
	}
	if len(result) != 0 {
		return result
	}
	return nil // 未找到
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
