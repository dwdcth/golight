package encrypt

import (
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/forgoer/openssl"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
)

var decoder *encoding.Encoder

func init() {
	decoder = unicode.UTF8.NewEncoder()
}

//加密字符串
func Encrypt(src string, key string) (string, error) {
	if src == "" {
		return "", nil
	}

	if key == "" {
		return "", errors.New("加密key不能为空")
	}

	plainText, err := decoder.String(src)
	if err != nil {
		return "", fmt.Errorf("源字符串UTF8编码转换失败 %s", err)
	}

	keyBytes, iv, err := getKeyIv(key)
	if err != nil {
		return "", fmt.Errorf("加密key UTF8编码转换失败 %s", err)
	}

	r, err := openssl.Des3CBCEncrypt([]byte(plainText), keyBytes, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("加密失败 %s", err)
	}

	return base64.StdEncoding.EncodeToString(r), nil
}

//解密字符串
func Decrypt(enStr string, key string) (string, error) {
	if enStr == "" {
		return "", nil
	}

	if key == "" {
		return "", errors.New("加密key不能为空")
	}

	keyBytes, iv, err := getKeyIv(key)
	if err != nil {
		return "", fmt.Errorf("加密key UTF8编码转换失败 %s", err)
	}

	srcBytes, err := base64.StdEncoding.DecodeString(enStr)
	if err != nil {
		return "", errors.New("待解密字符串UTF8编码转换失败")
	}

	r, err := openssl.Des3CBCDecrypt(srcBytes, keyBytes, iv, openssl.PKCS7_PADDING)
	if err != nil {
		return "", fmt.Errorf("加密失败 %s", err)
	}

	plainBytes, err := decoder.Bytes(r)
	if err != nil {
		return "", fmt.Errorf("待解密字符串UTF8编码转换失败 %s", err)
	}

	return string(plainBytes), nil
}

//获取keyBytes,iv
func getKeyIv(key string) (keyBytes []byte, iv []byte, err error) {
	keyUtf8, err := decoder.String(key)
	if err != nil {
		return nil, nil, fmt.Errorf("加密key UTF8编码转换失败 %s", err)
	}
	keyBytes = []byte(keyUtf8)
	iv = []byte(keyBytes[0:8])

	return keyBytes, iv, nil
}
