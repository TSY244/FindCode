package util

import (
	"context"
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	GitClonePrefix = "git-clone-"
)

func CloneRepository(ctx context.Context, url string) (string, error) {
	tempDir, err := os.MkdirTemp("./", GitClonePrefix)
	if err != nil {
		return "", err
	}
	sshAuth, err := createSSHAuth()
	if err != nil {
		return "", err
	}

	if err := cloneRepository(url, tempDir, sshAuth); err != nil {
		return "", err
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

// 创建 SSH 认证（支持带密码的密钥）
func createSSHAuth() (transport.AuthMethod, error) {
	// 1. 尝试从环境变量获取密钥路径
	keyPath := os.Getenv("SSH_KEY_PATH")
	if keyPath == "" {
		// 2. 使用默认路径 ~/.ssh/id_rsa
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户目录失败: %w", err)
		}
		keyPath = filepath.Join(home, ".ssh", "id_rsa")
	}

	// 3. 检查密钥文件是否存在
	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("SSH 密钥不存在: %s", keyPath)
	}

	// 4. 从环境变量获取密钥密码（没有密码则为空）
	keyPassword := os.Getenv("SSH_KEY_PASSWORD")

	// 5. 创建认证对象
	publicKeys, err := ssh.NewPublicKeysFromFile("git", keyPath, keyPassword)
	if err != nil {
		return nil, fmt.Errorf("加载 SSH 密钥失败: %w", err)
	}

	// 6. 忽略主机密钥验证（可选，根据环境需要）
	//publicKeys.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	return publicKeys, nil
}

func cloneRepository(url, path string, auth transport.AuthMethod) error {
	// 特殊处理：当使用 SSH 且 URL 是 HTTPS 时自动转换
	if strings.HasPrefix(url, "https://") && auth != nil {
		if _, ok := auth.(*ssh.PublicKeys); ok {
			url = strings.Replace(url, "https://", "git@", 1)
			url = strings.Replace(url, "/", ":", 1)
		}
	}

	_, err := git.PlainClone(path, false, &git.CloneOptions{
		URL:      url,
		Auth:     auth,
		Progress: os.Stdout,
	})
	return err
}
