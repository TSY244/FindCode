package utils

import (
	"errors"
	"fmt"
	"strings"
)

func extractJSON(input string) (string, error) {
	inString := false
	escape := false
	level := 0
	startIndex := -1
	endIndex := -1

	for i, c := range input {
		if escape {
			escape = false
			continue
		}

		switch c {
		case '\\':
			if inString {
				escape = true
			}
		case '"':
			inString = !inString
		case '{':
			if !inString {
				level++
				if level == 1 {
					startIndex = i
				}
			}
		case '}':
			if !inString && level > 0 {
				level--
				if level == 0 {
					endIndex = i
				}
			}
		}
	}

	if startIndex == -1 {
		return "", errors.New("no JSON object found")
	}
	if level != 0 {
		return "", errors.New("unmatched braces in JSON")
	}
	if endIndex == -1 {
		return "", errors.New("invalid JSON structure")
	}

	return input[startIndex : endIndex+1], nil
}

func ExtractJSON(input string) string {

	// 清理输入中的字面转义符（根据实际输入格式可能需要）
	cleaned := strings.ReplaceAll(input, `\n`, "\n")
	cleaned = strings.ReplaceAll(cleaned, `\"`, `"`)

	jsonStr, err := extractJSON(cleaned)
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	return jsonStr
}
