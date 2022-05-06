package slice

/**
 * 判断某个数值是否在切片中
 */
func IsValueInList(value int, list []int) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}
func IsStringValueInList(value string, list []string) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func SliceUnique(s []string) []string {
	temp := make([]string, 0)
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range s {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			temp = append(temp, e)
		}
	}
	return temp
}
