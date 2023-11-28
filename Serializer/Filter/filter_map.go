package Filter

import (
	"reflect"
)

func filterByGroups_map[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
		elemType = elemType.Elem()
	}

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
