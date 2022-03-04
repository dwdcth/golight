package validrule

import (
	"errors"
	"strconv"
	"strings"

	"github.com/elfincafe/mbstring"
)

// MbLen 字符串长度
// rule: 规则字符串
// value: 待验证的值
// message: 验证失败的提示信息
// data: 参数表示校验时传递的所有参数，例如校验的是一个map或者struct时，往往在联合校验时有用
func MbLen(rule string, value interface{}, message string, data interface{}) error {
	var src = value.(string)
	rule = strings.ReplaceAll(rule, "mbLen:", "")
	rules := strings.Split(rule, ",")
	srcLen := mbstring.Length(src)
	if message == "" {
		message = "mbLen valid failed"
	}
	min, _ := strconv.Atoi(rules[0])

	if len(rules) == 1 {
		if srcLen > min {
			return errors.New(message)
		}
		return nil
	} else {
		max, _ := strconv.Atoi(rules[1])
		if srcLen < min || srcLen > max {
			return errors.New(message)
		}
		return nil
	}

}
