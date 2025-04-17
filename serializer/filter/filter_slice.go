package filter

import (
	"reflect"
)

func filterByGroupsSlice[T any](obj T, groups ...string) T {
	value := DereferenceValueIfPointer(reflect.ValueOf(obj))

	if value.Len() == 0 {
		return obj
	}

	firstElem := value.Index(0).Interface()

	filteredFirstElem := FilterByGroups(firstElem, groups...)

	newSlice := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(filteredFirstElem)), 0, value.Len())

	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		filteredElem := FilterByGroups(elem.Interface(), groups...)
		newSlice = reflect.Append(newSlice, reflect.ValueOf(filteredElem))
	}
	return newSlice.Interface().(T)
}
