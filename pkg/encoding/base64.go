package encoding

import "encoding/base64"

// Base64Encode Base64 encode the input string
func Base64Encode(input string) string {
	return base64.StdEncoding.EncodeToString([]byte(input))
}

// Base64Decode Base64 decode the input string
func Base64Decode(input string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(input)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
