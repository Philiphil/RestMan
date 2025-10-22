package filter

import (
	"reflect"
	"strings"
)

// IsFieldIncluded checks if a struct field should be included based on its group tags and the provided groups.
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

// IsStruct checks if a type is a struct, dereferencing pointers if necessary.
func IsStruct(t reflect.Type) bool {
	return DereferenceTypeIfPointer(t).Kind() == reflect.Struct
}

// IsList checks if a type is a slice or array, dereferencing pointers if necessary.
func IsList(t reflect.Type) bool {
	return DereferenceTypeIfPointer(t).Kind() == reflect.Slice || DereferenceTypeIfPointer(t).Kind() == reflect.Array
}

// IsMap checks if a type is a map, dereferencing pointers if necessary.
func IsMap(t reflect.Type) bool {
	return DereferenceTypeIfPointer(t).Kind() == reflect.Map
}

func isAnonymous(field reflect.StructField) bool {
	return field.Anonymous
}

// DereferenceValueIfPointer recursively dereferences a value if it is a pointer.
func DereferenceValueIfPointer(value reflect.Value) reflect.Value {
	if value.Kind() == reflect.Ptr {
		return DereferenceValueIfPointer(value.Elem())
	}
	return value
}

// DereferenceTypeIfPointer recursively dereferences a type if it is a pointer.
func DereferenceTypeIfPointer(t reflect.Type) reflect.Type {
	if t.Kind() == reflect.Ptr {
		return DereferenceTypeIfPointer(t.Elem())
	}
	return t
}
