package filter

import (
	"fmt"
	"reflect"
)

func AssignFieldValue(field reflect.StructField, destValue reflect.Value, srcValue reflect.Value) error {
	if srcValue.IsZero() {
		return nil
	}
	if field.Type == srcValue.Type() || field.Type.AssignableTo(srcValue.Type()) {
		destValue.Set(srcValue)
	} else if field.Type.Kind() == reflect.Ptr && DereferenceTypeIfPointer(field.Type) == DereferenceTypeIfPointer(srcValue.Type()) {
		newPtr := reflect.New(field.Type.Elem())
		newPtr.Elem().Set(srcValue)
		destValue.Set(newPtr)
	} else if IsStruct(srcValue.Type()) {
		destFieldType := DereferenceTypeIfPointer(destValue.Type())
		if destFieldType.Kind() == reflect.Struct {
			destValueConverted := reflect.New(destFieldType).Elem()
			for i := 0; i < destFieldType.NumField(); i++ {
				destField := destFieldType.Field(i)
				srcFieldValue := DereferenceValueIfPointer(srcValue)

				srcFieldValue = srcFieldValue.FieldByName(destField.Name)
				if !srcFieldValue.IsValid() { //there no way this should happen but the linter yells at me
					continue
				}
				destFieldValue := destValueConverted.Field(i)

				if IsStruct(srcFieldValue.Type()) && destFieldValue.Kind() == reflect.Struct {
					AssignFieldValue(destField, destFieldValue, srcFieldValue)
				} else {
					AssignFieldValue(destField, destFieldValue, srcFieldValue)
				}
			}

			destValue.Set(destValueConverted)
		} else {
			//mostly for date objects
			destValue.Set(srcValue.Convert(destFieldType))
		}
	} else {
		if !srcValue.Type().ConvertibleTo(destValue.Type()) {
			return fmt.Errorf("type conversion not supported from %v to %v", srcValue.Type(), destValue.Type())
		}
		destValue.Set(srcValue.Convert(destValue.Type()))
	}
	return nil
}
