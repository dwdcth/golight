package encrypt

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"sort"
)

// md5参数加签，返回加签后的md5字符串
// paramMap 待校验参数
// appKey 应用的appkey的名称
// appVal 应用的appkey的值
// secretKey 应用的secretkey的名称
// secretVal 应用的secretkey的值
func Md5ParamSign(paramMap map[string]interface{}, appKey string, appVal string, secretKey, secretVal string) string {
	var keys []string
	for k := range paramMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var signStr string
	for _, k := range keys {
		signStr += k + "=" + fmt.Sprintf("%v", paramMap[k]) + "&"
	}
	signStr += fmt.Sprintf("%s=%s&", appKey, appVal)
	signStr += fmt.Sprintf("%s=%s", secretKey, secretVal)
	return fmt.Sprintf("%x", md5.Sum([]byte(signStr)))
}

func Md5FileSum(filepath string) (string, error) {
	f, err := os.Open(filepath)
	if err != nil {
		str1 := "Open err"
		return str1, err
	}
	defer f.Close()

	body, err := ioutil.ReadAll(f)
	if err != nil {
		str2 := "ioutil.ReadAll"
		return str2, err
	}
	md5 := fmt.Sprintf("%x", md5.Sum(body))
	runtime.GC()
	return md5, nil
}
