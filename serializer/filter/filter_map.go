package filter

import (
	"reflect"
)

func filterByGroupsMap[T any](obj T, groups ...string) T {
	value := DereferenceValueIfPointer(reflect.ValueOf(obj))
	elemType := DereferenceTypeIfPointer(value.Type())

	mapType := reflect.MapOf(elemType.Key(), elemType.Elem())
	mapValue := reflect.MakeMap(mapType)

	iter := value.MapRange()

	for iter.Next() {
		key := iter.Key()
		val := iter.Value()
		filteredVal := FilterByGroups(val.Interface(), groups...)

		filteredValValue, ok := filteredVal.(reflect.Value)
		if !ok {
			filteredValValue = reflect.ValueOf(filteredVal)
		}

		//Somehow unspecified struct cannot be assigned to the specified version of it
		if filteredValValue.Type().AssignableTo(elemType.Elem()) {
			mapValue.SetMapIndex(key, filteredValValue)
		} else if filteredValValue.Kind() == reflect.Struct && elemType.Elem().Kind() == reflect.Struct {
			// Both are structs but different types (filtered vs original)
			// Need to copy field by field matching by name
			destType := elemType.Elem()
			destValue := reflect.New(destType).Elem()

			filteredValValue = DereferenceValueIfPointer(filteredValValue)

			for i := 0; i < destType.NumField(); i++ {
				destField := destType.Field(i)
				destFieldValue := destValue.Field(i)

				// Try to find matching field in filtered struct
				srcFieldValue := filteredValValue.FieldByName(destField.Name)
				if srcFieldValue.IsValid() && !srcFieldValue.IsZero() {
					assignFieldValue(destField, destFieldValue, srcFieldValue)
				}
			}

			mapValue.SetMapIndex(key, destValue)
		} else {
			// Try direct conversion
			if filteredValValue.Type().ConvertibleTo(elemType.Elem()) {
				mapValue.SetMapIndex(key, filteredValValue.Convert(elemType.Elem()))
			} else {
				// Last resort - try to assign via assignFieldValue with first field
				destType := elemType.Elem()
				destValue := reflect.New(destType).Elem()
				if destType.NumField() > 0 {
					assignFieldValue(destType.Field(0), destValue.Field(0), filteredValValue)
				}
				mapValue.SetMapIndex(key, destValue)
			}
		}
	}

	return mapValue.Interface().(T)
}
