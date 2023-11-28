package Filter

import (
	"fmt"
	"reflect"
)

func assignFieldValue(field reflect.StructField, destValue reflect.Value, srcValue reflect.Value) {
	if srcValue.IsZero() {
		return
	}
	if field.Type == srcValue.Type() {
		destValue.Set(srcValue)
	} else if field.Type.AssignableTo(srcValue.Type()) {
		destValue.Set(srcValue)
	} else if field.Type.Kind() == reflect.Ptr && field.Type.Elem() == srcValue.Type() {
		newPtr := reflect.New(field.Type.Elem())
		newPtr.Elem().Set(srcValue)
		destValue.Set(newPtr)
	} else if isStruct(srcValue.Type()) {
		destFieldType := destValue.Type()
		if destFieldType.Kind() == reflect.Ptr {
			destFieldType = destFieldType.Elem()
		}

		if destFieldType.Kind() == reflect.Struct {
			destValueConverted := reflect.New(destFieldType).Elem()

			for i := 0; i < destFieldType.NumField(); i++ {
				destField := destFieldType.Field(i)
				srcFieldValue := srcValue
				if srcFieldValue.Type().Kind() == reflect.Ptr {
					srcFieldValue = srcFieldValue.Elem()
				}

				srcFieldValue = srcFieldValue.FieldByName(destField.Name)
				destFieldValue := destValueConverted.Field(i)

				if isStruct(srcFieldValue.Type()) && destFieldValue.Kind() == reflect.Struct {
					assignFieldValue(destField, destFieldValue, srcFieldValue)
				} else {
					assignFieldValue(destField, destFieldValue, srcFieldValue)
				}
			}

			destValue.Set(destValueConverted)
		} else {
			destValue.Set(srcValue.Convert(destFieldType))
		}
	} else {
		if !srcValue.Type().ConvertibleTo(destValue.Type()) {
			fmt.Printf("Type conversion not supported from %v to %v\n", srcValue.Type(), destValue.Type())
			panic("!")
		}
		destValue.Set(srcValue.Convert(destValue.Type()))
	}
}
