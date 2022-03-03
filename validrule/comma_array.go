package validrule

import (
	"fmt"
	"regexp"
	"strings"
)

var commaIntegerArrayReg = regexp.MustCompile(`^\d+(,\d+)*$`)

// RuleCommaArray 整数都好分割数组验证
// rule: 规则字符串
// value: 待验证的值
// message: 验证失败的提示信息
// data: 参数表示校验时传递的所有参数，例如校验的是一个map或者struct时，往往在联合校验时有用
func RuleCommaArray(rule string, value interface{}, message string, data interface{}) error {
	var ruleParam string
	if ruleParams := strings.Split(rule, ":"); len(ruleParams) == 2 {
		ruleParam = ruleParams[1]
	} else {
		return fmt.Errorf("%s", "comma-array rule param is error")
	}
	switch ruleParam {
	case "integer":
		if !commaIntegerArrayReg.MatchString(value.(string)) {
			return fmt.Errorf("%s %s", message, "must be a comma-separated list of integers")
		}
	}
	return nil
}
