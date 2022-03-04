package encode

import (
	"fmt"
	"io/ioutil"
	"os"

	qrencode "github.com/skip2/go-qrcode"
	qrdecode "github.com/tuotoo/qrcode"
)

func QrEncode(content string, level qrencode.RecoveryLevel, size int) ([]byte, error) {
	var png []byte
	png, err := qrencode.Encode(content, level, size)
	if err == nil {
		return png, nil
	}
	return nil, err
}

func QrEncodeFile(content string, filename string, level qrencode.RecoveryLevel, size int) error {
	png, err := QrEncode(content, level, size)
	if err == nil {
		return ioutil.WriteFile(filename, png, os.FileMode(0644))
	}
	return err
}

func QrDecode(filePath string) string {
	fi, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	defer fi.Close()
	qr, err := qrdecode.Decode(fi)
	if err == nil {
		return qr.Content
	}
	return ""
}
