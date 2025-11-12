package filter

import (
	"reflect"
)

func filterByGroupsStruct[T any](obj T, groups []string, sparseFields map[string][]string, currentPath string) T {
	value := reflect.ValueOf(obj)
	elemType := value.Type()

	if elemType.Kind() == reflect.Pointer {
		if value.IsNil() {
			return obj
		}
		elemType = DereferenceTypeIfPointer(elemType)
		value = DereferenceValueIfPointer(value)
	}

	// Disable cache when sparse fields are specified to avoid incorrect filtering
	useCacheEnabled := len(sparseFields) == 0

	if useCacheEnabled {
		if cachedEntry, ok := globalCache.Get(elemType, groups); ok {
			newValue := reflect.New(cachedEntry.filteredType).Elem()
			if len(cachedEntry.fieldMappings) > 0 {
				populateStructFieldsFast(newValue, value, cachedEntry.fieldMappings)
			} else {
				populateStructFields(newValue, value, cachedEntry.filteredType)
			}
			// No sparse fields filtering needed since cache doesn't support it
			return newValue.Interface().(T)
		}
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
					nestedPath := field.Name
					if currentPath != "" {
						nestedPath = currentPath + "." + field.Name
					}
					filteredElem := filterByGroupsWithPath(fieldValue.Interface(), groups, sparseFields, nestedPath)
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

	// Only cache if sparse fields are not being used
	if useCacheEnabled {
		entry := &typeCacheEntry{
			filteredType:  newStructType,
			fieldMappings: mappings,
		}
		globalCache.Set(elemType, groups, entry)
	}

	newValue := reflect.New(newStructType).Elem()

	if originalTypeName != "" && len(newFields) > 0 {
	}

	if len(mappings) > 0 {
		populateStructFieldsFast(newValue, value, mappings)
	} else {
		populateStructFields(newValue, value, newStructType)
	}

	// Apply sparse fields filtering if specified
	if len(sparseFields) > 0 {
		filteredValue := applySparseFieldsFilter(newValue, sparseFields, currentPath)
		return filteredValue.Interface().(T)
	}

	return newValue.Interface().(T)
}

// applySparseFieldsFilter removes fields not requested in sparse fields.
// Fields are actually removed from the struct, not set to zero/nil.
func applySparseFieldsFilter(value reflect.Value, sparseFields map[string][]string, currentPath string) reflect.Value {
	requestedFields, hasSparseFields := sparseFields[currentPath]
	if !hasSparseFields {
		return value
	}

	// Build set of requested field names for O(1) lookup
	requestedSet := make(map[string]bool)
	for _, field := range requestedFields {
		requestedSet[field] = true
	}

	// Build new struct type with only requested fields
	var newFields []reflect.StructField
	var fieldMappings []fieldMapping
	valueType := value.Type()

	srcIdx := 0
	destIdx := 0
	for i := 0; i < valueType.NumField(); i++ {
		field := valueType.Field(i)

		// Check if field is requested (use JSON tag name if available)
		fieldName := field.Name
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			fieldName = jsonTag
		}

		if requestedSet[fieldName] || requestedSet[field.Name] {
			newFields = append(newFields, field)
			fieldMappings = append(fieldMappings, fieldMapping{srcIndex: srcIdx, destIndex: destIdx})
			destIdx++
		}
		srcIdx++
	}

	// If no fields match, return empty struct of same type
	if len(newFields) == 0 {
		return reflect.Zero(value.Type())
	}

	// Create new struct with only requested fields
	newStructType := reflect.StructOf(newFields)
	newValue := reflect.New(newStructType).Elem()

	// Copy field values
	for _, mapping := range fieldMappings {
		srcField := value.Field(mapping.srcIndex)
		destField := newValue.Field(mapping.destIndex)
		destField.Set(srcField)
	}

	return newValue
}

func populateStructFieldsFast(newValue, value reflect.Value, mappings []fieldMapping) {
	for _, mapping := range mappings {
		srcField := value.Field(mapping.srcIndex)
		destField := newValue.Field(mapping.destIndex)
		destFieldType := newValue.Type().Field(mapping.destIndex)
		AssignFieldValue(destFieldType, destField, srcField)
	}
}

func populateStructFields(newValue, value reflect.Value, structType reflect.Type) {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		fieldName := field.Name
		fieldValue := value.FieldByName(fieldName)
		newFieldValue := newValue.Field(i)
		AssignFieldValue(field, newFieldValue, fieldValue)
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
