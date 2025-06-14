package util

import (
	"ScanIDOR/pkg/logger"
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"regexp"
	"strings"
)

func CloneRepository(ctx context.Context, url string) (string, error) {
	tempDir, err := os.MkdirTemp("./", "git-clone-")
	if err != nil {
		return "", err
	}
	_, err = git.PlainCloneContext(ctx, tempDir, false, &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	})
	if err != nil {
		os.RemoveAll(tempDir) // 尽力清理
		// 5. 包装错误，提供上下文
		logger.Errorf("从 URL '%s' 克隆仓库失败: %v", url, err)
		return "", fmt.Errorf("从 URL '%s' 克隆仓库失败: %w", url, err)
	}
	return tempDir, nil
}

func CheckGitUrl(url string) bool {
	// 移除可能的.git后缀以便统一处理
	url = strings.TrimSuffix(url, ".git")

	// 改进后的正则表达式模式，支持：
	// 1. HTTP/HTTPS/GIT/SSH 协议格式
	// 2. SSH 的 SCP 格式 (user@host:path)
	// 3. 大写字母路径
	pattern := `^(?:https?|git|ssh)://(?:[a-zA-Z0-9._-]+@)?[a-zA-Z0-9.-]+(?::\d+)?(?:/[\w.-]+)+$|^(?:[\w.-]+@)?[\w.-]+:[\w./-]+$`

	// 编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false
	}

	return re.MatchString(url)
}
