package filter

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

func IsStruct(t reflect.Type) bool {
	return DereferenceTypeIfPointer(t).Kind() == reflect.Struct
}

func IsList(t reflect.Type) bool {
	return DereferenceTypeIfPointer(t).Kind() == reflect.Slice || DereferenceTypeIfPointer(t).Kind() == reflect.Array
}

func IsMap(t reflect.Type) bool {
	return DereferenceTypeIfPointer(t).Kind() == reflect.Map
}

func isAnonymous(field reflect.StructField) bool {
	return field.Anonymous
}

func DereferenceValueIfPointer(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return DereferenceValueIfPointer(value.Elem())
	}
	return value
}

func DereferenceTypeIfPointer(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return DereferenceTypeIfPointer(t.Elem())
	}
	return t
}
