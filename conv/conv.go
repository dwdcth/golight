package conv

import (
	"reflect"
	"strings"
)

//struct转map[string]string
func StructToMapStr(in interface{}, tag string) map[string]string {
	//dst :=make(map[string]interface{})
	//mergo.Map(&dst,source)
	//return dst

	out := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		//panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
		return nil
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(tag)
		if tagv != "" && tagv != "-" {
			if strings.Contains(tagv, "omitempty") && v.Field(i).IsNil() {
				continue
			}
			out[tagv] = v.Field(i).String()
		}
	}
	return out
}

// 有类型的slice 转为 interface slice
func ToInterfaceSlice(slice interface{}) []interface{} {
	s := reflect.ValueOf(slice)
	if s.Kind() != reflect.Slice {
		panic("InterfaceSlice() given a non-slice type")
	}

	ret := make([]interface{}, s.Len())

	for i := 0; i < s.Len(); i++ {
		ret[i] = s.Index(i).Interface()
	}

	return ret
}

//结构体slice转interface{} slice
// isPtr 为false slice里的是结构体, true 的时候是结构体指针
func StructSliceToInterfaceSlice(nodes interface{}, isPtr bool) []interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(nodes))
	switch sliceValue.Kind() {
	case reflect.Slice:
		res := make([]interface{}, 0)
		s := reflect.ValueOf(nodes)
		for i := 0; i < s.Len(); i++ {
			item := s.Index(i)

			if isPtr == true && item.Kind() != reflect.Ptr {
				panic("nodes must be  ptr slice")
			}
			element := item.Interface()
			res = append(res, element)
		}
		return res
	default:
		return nil
	}
}

//struct转map
func StructToMap(in interface{}, tag string) map[string]interface{} {
	//dst :=make(map[string]interface{})
	//mergo.Map(&dst,source)
	//return dst

	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// we only accept structs
	if v.Kind() != reflect.Struct {
		//panic(fmt.Errorf("ToMap only accepts structs; got %T", v))
		return nil
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		// gets us a StructField
		fi := typ.Field(i)
		tagv := fi.Tag.Get(tag)
		if tagv != "" && tagv != "-" {
			if strings.Contains(tagv, "omitempty") && v.Field(i).IsNil() {
				continue
			}
			out[tagv] = v.Field(i).Interface()
		}
	}
	return out

}

//struct slice 转map  slice
func StructSliceToMapSlice(source interface{}) []map[string]interface{} {
	sliceValue := reflect.Indirect(reflect.ValueOf(source))
	switch sliceValue.Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(source)
		var res = make([]map[string]interface{}, 0)
		for i := 0; i < s.Len(); i++ {
			element := s.Index(i).Interface()
			res = append(res, StructToMap(element, "json"))
		}
		return res
	default:
		return nil
	}
}

func GetElemType(a interface{}) reflect.Type {
	for t := reflect.TypeOf(a); ; {
		switch t.Kind() {
		case reflect.Ptr, reflect.Slice:
			t = t.Elem()
		default:
			return t
		}
	}
}
