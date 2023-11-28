package Filter

import (
	"reflect"
)

func filterByGroups_struct[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if elemType.Kind() == reflect.Ptr {
		if value.IsNil() {
			return obj
		}
		elemType = elemType.Elem()
		value = value.Elem()
	}

	var newFields []reflect.StructField

	if value.IsValid() {
		for i := 0; i < value.NumField(); i++ {
			field := elemType.Field(i)
			if isFieldExported(field) && IsFieldIncluded(field, groups) {
				fieldValue := value.Field(i)

				if isStruct(field.Type) && !isAnonymous(field) {
					filteredElem := FilterByGroups(fieldValue.Interface(), groups...)
					newFields = append(newFields, reflect.StructField{
						Name: field.Name,
						Type: reflect.TypeOf(filteredElem),
						Tag:  field.Tag,
					})
				} else {
					newFields = append(newFields, field)
				}
			}
		}
		anonymousFields := filterAnonymousFields(value, groups...)
		newFields = append(newFields, anonymousFields...)
	}

	newStructType := reflect.StructOf(newFields)
	newValue := reflect.New(newStructType).Elem()

	for i, field := range newFields {
		fieldName := field.Name
		fieldValue := value.FieldByName(fieldName)
		newFieldValue := newValue.Field(i)
		assignFieldValue(field, newFieldValue, fieldValue)
	}

	return newValue.Interface().(T)
}

func filterAnonymousFields(value reflect.Value, groups ...string) []reflect.StructField {
	var anonymousFields []reflect.StructField

	if value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		if isAnonymous(field) {
			var fieldType reflect.Type
			if fieldValue.Kind() == reflect.Ptr {
				fieldType = fieldValue.Type().Elem()
			} else {
				fieldType = fieldValue.Type()
			}

			for j := 0; j < fieldType.NumField(); j++ {
				anonymousField := fieldType.Field(j)
				if isFieldExported(anonymousField) && IsFieldIncluded(anonymousField, groups) {
					anonymousFields = append(anonymousFields, reflect.StructField{
						Name: anonymousField.Name,
						Type: anonymousField.Type,
						Tag:  anonymousField.Tag,
					})
				}
			}
		}
	}

	return anonymousFields
}
