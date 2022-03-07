package encrypt

import (
	"fmt"
	"testing"
)

func TestMd5ParamSign(t *testing.T) {
	appSecret := "VABoXAxwfdwcqfVLrvgrheIaerKafjtFtQXNrZuAYbSDpaCyYZqJejPVfgnnkgIXlkCIdSVVnJwrWjZqVHdGSoEsbhqPeTKZBfodCkhnKojctTXTtBxijvlzxgQednxg"
	appkey := "eMamLGwrNaVGOFOlTCKnOEuivLekNtQupfOrXjkynNFMVbMi"
	sign := Md5ParamSign(map[string]interface{}{
		"IsDisable": "1",
		"Uuid":      "24743",
	}, "AppKey", appkey, "AppSecret", appSecret)
	fmt.Println(sign)
}
