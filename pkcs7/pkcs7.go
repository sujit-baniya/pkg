package pkcs7

import (
	"bytes"
	"errors"
)

func isPow2(num int) bool {
	return num != 0 && ((num & (num - 1)) == 0)
}

// Pad applies pkcs7 padding
func Pad(buf []byte, size int) ([]byte, error) {
	if !isPow2(size) {
		return nil, errors.New("size is not power of 2")
	}
	bufLen := len(buf)
	padLen := size - bufLen%size
	padText := bytes.Repeat([]byte{byte(padLen)}, padLen)
	return append(buf, padText...), nil
}

// Unpad removes pkcs7 padding
func Unpad(padded []byte, size int) ([]byte, error) {
	if !isPow2(size) {
		return nil, errors.New("size is not power of 2")
	}
	if len(padded)%size != 0 {
		return nil, errors.New("padded value wasn't in correct size")
	}
	paddedLen := len(padded)
	padLen := int(padded[paddedLen-1])
	bufLen := paddedLen - padLen
	return padded[:bufLen], nil
}
