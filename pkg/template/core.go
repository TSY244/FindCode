package template

import (
	"regexp"
	"strings"
)

// 用于替换{{}}

type Template struct {
	Template string
	Params   map[string]string
	Result   string
	Keys     []string // 存储模版中的key
	re       *regexp.Regexp
}

func NewTemplate(template string, params map[string]string) *Template {
	return &Template{
		Template: template,
		Params:   params,
	}
}

func (t *Template) Load() {
	t.re = regexp.MustCompile(`{{(\w+)}}`)
	temp := t.re.FindAllString(t.Template, -1)
	// 处理keys
	for _, v := range temp {
		t.Keys = append(t.Keys, v[2:len(v)-2])
	}
}

func (t *Template) Replace() (string, error) {
	// 编译支持点号的正则表达式
	re := regexp.MustCompile(`\{\{\s*([\w\.]+)\s*\}\}`)

	result := re.ReplaceAllStringFunc(t.Template, func(match string) string {
		// 提取纯键名（去除 {{ }} 和空格）
		key := re.ReplaceAllString(match, "$1")
		key = strings.TrimSpace(key)

		// 调试输出
		// fmt.Printf("匹配到占位符: %s, 提取键名: %s\n", match, key)

		if value, exists := t.Params[key]; exists {
			return value
		}

		// 未找到时返回原始占位符（或抛出错误）
		return match
	})

	return result, nil
}
