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
		} else {
			destType := elemType.Elem()
			destValue := reflect.New(destType).Elem()

			assignFieldValue(destType.Field(0), destValue.Field(0), filteredValValue)

			mapValue.SetMapIndex(key, destValue)
		}
	}

	return mapValue.Interface().(T)
}
