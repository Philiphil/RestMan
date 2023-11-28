package Filter

import (
	"reflect"
	"strings"
)

func IsFieldIncluded(field reflect.StructField, groups []string) bool {
	if len(groups) == 0 {
		return true //No filtration then
	}

	tag := field.Tag.Get("group")
	if tag == "" {
		return false
	}

	groupList := strings.Split(tag, ",")
	for _, group := range groups {
		for _, g := range groupList {
			if group == g {
				return true
			}
		}
	}

	return false
}

func isFieldExported(field reflect.StructField) bool {
	return field.PkgPath == ""
}

func isStruct(t reflect.Type) bool {
	return t.Kind() == reflect.Struct || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Struct)
}

func isList(t reflect.Type) bool {
	return t.Kind() == reflect.Slice || t.Kind() == reflect.Array || (t.Kind() == reflect.Ptr && (t.Elem().Kind() == reflect.Slice || t.Elem().Kind() == reflect.Array))
}
func isMap(t reflect.Type) bool {
	return t.Kind() == reflect.Map || (t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Map)
}

func isAnonymous(field reflect.StructField) bool {
	return field.Anonymous
}
