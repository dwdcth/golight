package filex

import (
	"errors"
	"io"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

// IsAllowFileExt 文件扩展名判断
// fileReader:文件读取器
// fileName:文件表面上的文件名
// allowFileExts:允许的文件扩展名
func IsAllowFileExt(fileReader io.Reader, fileName string, allowFileExts []string) (bool, error) {
	names := strings.Split(fileName, ".")
	if len(names) < 2 {
		return false, errors.New("文件名不合法")
	}
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
	ext := fileDetect.Extension()
	found := false
	for _, e := range allowFileExts {
		if e == ext {
			found = true
			break
		}
	}
	if !strings.Contains(ext, names[len(names)-1]) {
		found = false
	}
	return found, nil
}
