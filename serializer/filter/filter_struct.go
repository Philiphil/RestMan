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

	if cachedEntry, ok := globalCache.Get(elemType, groups); ok {
		newValue := reflect.New(cachedEntry.filteredType).Elem()
		if len(cachedEntry.fieldMappings) > 0 {
			populateStructFieldsFast(newValue, value, cachedEntry.fieldMappings)
		} else {
			populateStructFields(newValue, value, cachedEntry.filteredType, groups)
		}
		return newValue.Interface().(T)
	}

	var newFields []reflect.StructField
	var mappings []fieldMapping
	originalTypeName := elemType.Name()

	if value.IsValid() {
		destIdx := 0
		for i := range value.NumField() {
			field := elemType.Field(i)
			if isFieldExported(field) && IsFieldIncluded(field, groups) && !isAnonymous(field) {
				fieldValue := value.Field(i)

				if IsStruct(field.Type) {
					filteredElem := FilterByGroups(fieldValue.Interface(), groups...)
					newField := reflect.StructField{
						Name: field.Name,
						Type: reflect.TypeOf(filteredElem),
						Tag:  field.Tag,
					}
					if newField.Tag.Get("xml") == "" {
						if newField.Tag == "" {
							newField.Tag = reflect.StructTag(`xml:"` + field.Name + `"`)
						} else {
							newField.Tag = reflect.StructTag(string(newField.Tag) + ` xml:"` + field.Name + `"`)
						}
					}
					newFields = append(newFields, newField)
					mappings = append(mappings, fieldMapping{srcIndex: i, destIndex: destIdx})
					destIdx++
				} else {
					newField := field
					if newField.Tag.Get("xml") == "" {
						if newField.Tag == "" {
							newField.Tag = reflect.StructTag(`xml:"` + field.Name + `"`)
						} else {
							newField.Tag = reflect.StructTag(string(newField.Tag) + ` xml:"` + field.Name + `"`)
						}
					}
					newFields = append(newFields, newField)
					mappings = append(mappings, fieldMapping{srcIndex: i, destIndex: destIdx})
					destIdx++
				}
			}
		}
		anonymousFields := filterAnonymousFields(value, groups...)
		if len(anonymousFields) > 0 {
			mappings = nil
		}
		newFields = append(newFields, anonymousFields...)
	}

	newStructType := reflect.StructOf(newFields)
	entry := &typeCacheEntry{
		filteredType: newStructType,
		fieldMappings: mappings,
	}
	globalCache.Set(elemType, groups, entry)
	newValue := reflect.New(newStructType).Elem()

	if originalTypeName != "" && len(newFields) > 0 {
	}

	if len(mappings) > 0 {
		populateStructFieldsFast(newValue, value, mappings)
	} else {
		populateStructFields(newValue, value, newStructType, groups)
	}

	return newValue.Interface().(T)
}

func populateStructFieldsFast(newValue, value reflect.Value, mappings []fieldMapping) {
	for _, mapping := range mappings {
		srcField := value.Field(mapping.srcIndex)
		destField := newValue.Field(mapping.destIndex)
		destFieldType := newValue.Type().Field(mapping.destIndex)
		assignFieldValue(destFieldType, destField, srcField)
	}
}

func populateStructFields(newValue, value reflect.Value, structType reflect.Type, groups []string) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name
		fieldValue := value.FieldByName(fieldName)
		newFieldValue := newValue.Field(i)
		assignFieldValue(field, newFieldValue, fieldValue)
	}
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
