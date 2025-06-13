package utils

import (
	"bytes"
	"compress/zlib"
	"io"
)

func Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	w.Close()
	return buf.Bytes(), nil
}

func Decompress(data []byte) ([]byte, error) {
	var out bytes.Buffer
	r, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	_, err = io.Copy(&out, r)
	if err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}
