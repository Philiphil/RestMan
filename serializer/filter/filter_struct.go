package filter

import (
	"reflect"
)

func filterByGroupsStruct[T any](obj T, groups ...string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if elemType.Kind() == reflect.Ptr {
		if value.IsNil() {
			return obj
		}
		elemType = DereferenceTypeIfPointer(elemType)
		value = DereferenceValueIfPointer(value)
	}

	var newFields []reflect.StructField
	originalTypeName := elemType.Name()

	if value.IsValid() {
		for i := range value.NumField() {
			field := elemType.Field(i)
			if isFieldExported(field) && IsFieldIncluded(field, groups) {
				fieldValue := value.Field(i)

				if IsStruct(field.Type) && !isAnonymous(field) {
					filteredElem := FilterByGroups(fieldValue.Interface(), groups...)
					newField := reflect.StructField{
						Name: field.Name,
						Type: reflect.TypeOf(filteredElem),
						Tag:  field.Tag,
					}
					// Add xml tag with original field name if not present
					if newField.Tag.Get("xml") == "" {
						if newField.Tag == "" {
							newField.Tag = reflect.StructTag(`xml:"` + field.Name + `"`)
						} else {
							newField.Tag = reflect.StructTag(string(newField.Tag) + ` xml:"` + field.Name + `"`)
						}
					}
					newFields = append(newFields, newField)
				} else {
					newField := field
					// Add xml tag if not present
					if newField.Tag.Get("xml") == "" {
						if newField.Tag == "" {
							newField.Tag = reflect.StructTag(`xml:"` + field.Name + `"`)
						} else {
							newField.Tag = reflect.StructTag(string(newField.Tag) + ` xml:"` + field.Name + `"`)
						}
					}
					newFields = append(newFields, newField)
				}
			}
		}
		anonymousFields := filterAnonymousFields(value, groups...)
		newFields = append(newFields, anonymousFields...)
	}

	newStructType := reflect.StructOf(newFields)
	newValue := reflect.New(newStructType).Elem()

	// Set XMLName if this was a named struct
	if originalTypeName != "" && len(newFields) > 0 {
		// Try to add XMLName field at the beginning for proper XML marshaling
		// Note: reflect.StructOf doesn't support adding XMLName after creation
	}

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

	value = DereferenceValueIfPointer(value)

	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		fieldValue := value.Field(i)

		if isAnonymous(field) {
			fieldType := DereferenceTypeIfPointer(fieldValue.Type())

			for j := 0; j < fieldType.NumField(); j++ {
				anonymousField := fieldType.Field(j)
				if isFieldExported(anonymousField) && IsFieldIncluded(anonymousField, groups) {
					newField := reflect.StructField{
						Name: anonymousField.Name,
						Type: anonymousField.Type,
						Tag:  anonymousField.Tag,
					}
					// Add xml tag if not present
					if newField.Tag.Get("xml") == "" {
						if newField.Tag == "" {
							newField.Tag = reflect.StructTag(`xml:"` + anonymousField.Name + `"`)
						} else {
							newField.Tag = reflect.StructTag(string(newField.Tag) + ` xml:"` + anonymousField.Name + `"`)
						}
					}
					anonymousFields = append(anonymousFields, newField)
				}
			}
		}
	}

	return anonymousFields
}
