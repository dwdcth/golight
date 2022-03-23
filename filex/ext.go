package filex

import (
	"errors"
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// 注意fileReader会更改位置，如果其它地方需要使用，需要Seek(0,0)
// IsAllowFileExt 文件扩展名判断
// fileReader:文件读取器
// fileName:文件表面上的文件名
// allowFileExts:允许的文件扩展名
func IsAllowFileExt(fileReader io.Reader, fileName string, allowFileExts []string) (bool, error) {
	names := strings.Split(fileName, ".")
	if len(names) < 2 {
		return false, errors.New("文件名不合法")
	}
	fileExt := names[len(names)-1]
	fileDetect, err := mimetype.DetectReader(fileReader)
	if err != nil {
		return false, err
	}
	if fileDetect == nil {
		return false, errors.New("fileReader error")
	}
	if allowFileExts == nil {
		return false, errors.New("allowFileExts error")
	}
	detectExt := fileDetect.Extension()
	found := false
	for _, e := range allowFileExts {
		if e == detectExt {
			found = true
			break
		}
	}
	if !strings.Contains(detectExt, fileExt) {
		found = false
		if detectExt == ".jpg" && fileExt == "jpeg" {
			found = true
		}
		if detectExt == ".txt" && fileExt == "csv" {
			found = true
		}
	}
	return found, nil
}
