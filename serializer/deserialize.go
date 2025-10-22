package serializer

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"reflect"
	"strings"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/serializer/filter"
)

// Deserialize converts a string representation into an object in the configured format.
func (s *Serializer) Deserialize(data string, obj any) error {
	if !isPointer(obj) {
		return fmt.Errorf("object must be pointer")
	}
	switch s.Format {
	case format.JSONLD:
		return json.Unmarshal([]byte(data), obj)
	case format.JSON:
		return json.Unmarshal([]byte(data), obj)
	case format.XML:
		return s.deserializeXML(data, obj)
	case format.CSV:
		return s.deserializeCSV(data, obj)
	default:
		return fmt.Errorf("unsupported format: %s", s.Format)
	}
}

func (s *Serializer) deserializeXML(data string, obj any) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("obj must be a pointer")
	}

	elemValue := value.Elem()

	// If it's a slice, we need special handling for <items><item>...</item></items> format
	if elemValue.Kind() == reflect.Slice {
		// Create a wrapper struct to match our XML structure
		sliceElemType := elemValue.Type().Elem()

		// Try to unmarshal into a temporary structure
		decoder := xml.NewDecoder(strings.NewReader(data))

		// Skip to first element
		for {
			token, err := decoder.Token()
			if err != nil {
				return err
			}

			if _, ok := token.(xml.StartElement); ok {
				// Found root element, now read items
				newSlice := reflect.MakeSlice(elemValue.Type(), 0, 0)

				for {
					token, err := decoder.Token()
					if err != nil {
						if err.Error() == "EOF" {
							break
						}
						return err
					}

					if itemStart, ok := token.(xml.StartElement); ok {
						if itemStart.Name.Local == "item" {
							// Decode this item
							newElem := reflect.New(sliceElemType).Interface()
							if err := decoder.DecodeElement(newElem, &itemStart); err != nil {
								return err
							}
							newSlice = reflect.Append(newSlice, reflect.ValueOf(newElem).Elem())
						}
					} else if _, ok := token.(xml.EndElement); ok {
						// End of root element
						break
					}
				}

				elemValue.Set(newSlice)
				return nil
			}
		}
	}

	// For non-slice types, use standard unmarshaling
	return xml.Unmarshal([]byte(data), obj)
}

// MergeObjects merges source object fields into target object, both must be pointers.
func (s *Serializer) MergeObjects(target any, source any) error {
	targetValue := reflect.ValueOf(target)
	sourceValue := reflect.ValueOf(source)

	if targetValue.Kind() != reflect.Ptr || sourceValue.Kind() != reflect.Ptr {
		return fmt.Errorf("both target and source must be pointers")
	}

	targetValue = targetValue.Elem()
	sourceValue = sourceValue.Elem()

	mergeFields(targetValue, sourceValue)

	return nil
}

func mergeFields(target reflect.Value, source reflect.Value) {
	source = filter.DereferenceValueIfPointer(source)

	//if target is nil or empty, lets create anew
	if (target.Kind() == reflect.Ptr || target.Kind() == reflect.Interface) && target.IsNil() {
		newTarget := reflect.New(source.Type())
		if target.Kind() == reflect.Ptr {
			target.Set(newTarget)
		}
		target = newTarget.Elem()
	} else if target.Kind() == reflect.Slice && target.IsNil() {
		target.Set(reflect.MakeSlice(source.Type(), 0, source.Len()))
	}

	target = filter.DereferenceValueIfPointer(target)

	if target.Kind() == reflect.Struct && source.Kind() == reflect.Struct {
		for i := 0; i < target.NumField(); i++ {
			targetField := target.Field(i)
			sourceField := source.Field(i)

			if shouldExclude(targetField) {
				continue
			}

			if targetField.CanSet() && !isEmpty(sourceField) {
				if targetField.Kind() == reflect.Struct && sourceField.Kind() == reflect.Struct {
					mergeFields(targetField, sourceField)
				} else {
					targetField.Set(sourceField)
				}
			}
		}
		return
	}

	if target.Kind() == reflect.Slice && source.Kind() == reflect.Slice {
		for i := 0; i < source.Len(); i++ {
			sourceElem := source.Index(i)
			if sourceElem.Kind() == reflect.Ptr || sourceElem.Kind() == reflect.Struct {
				mergedElem := reflect.New(sourceElem.Type()).Elem()
				if i < target.Len() {
					mergeFields(target.Index(i), sourceElem)
				} else {
					mergeFields(mergedElem, sourceElem)
					target.Set(reflect.Append(target, mergedElem))
				}
			} else {
				target.Set(reflect.Append(target, sourceElem))
			}
		}
	}
}

func shouldExclude(field reflect.Value) bool {
	fieldName := field.Type().Name()
	excludedFields := []string{"CreatedAt", "ModifiedAt", "DeletedAt"}

	for _, excluded := range excludedFields {
		if strings.EqualFold(fieldName, excluded) {
			return true
		}
	}
	return false
}

// DeserializeAndMerge deserializes data and merges it into the target object.
func (s *Serializer) DeserializeAndMerge(data string, target any) error {
	source := reflect.New(reflect.TypeOf(target).Elem()).Interface()

	if err := s.Deserialize(data, source); err != nil {
		return err
	}

	return s.MergeObjects(target, source)
}

// isEmpty checks if a value has the zero value of its type
func isEmpty(v reflect.Value) bool {
	zero := reflect.Zero(v.Type())
	return reflect.DeepEqual(v.Interface(), zero.Interface())
}

// isPointer checks if a value is a pointer
func isPointer(v any) bool {
	t := reflect.TypeOf(v)
	return t.Kind() == reflect.Ptr
}

func (s *Serializer) deserializeCSV(data string, obj any, groups ...string) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("CSV deserialization requires a non-nil pointer")
	}

	sliceValue := value.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("CSV deserialization requires a pointer to a slice, got %s", sliceValue.Kind())
	}

	reader := csv.NewReader(strings.NewReader(data))
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	if len(rows) == 0 {
		return nil
	}

	// First row is header
	header := rows[0]
	if len(header) == 0 {
		return fmt.Errorf("CSV header is empty")
	}

	// Get element type
	elemType := sliceValue.Type().Elem()
	isPtr := elemType.Kind() == reflect.Ptr
	if isPtr {
		elemType = elemType.Elem()
	}

	if elemType.Kind() != reflect.Struct {
		return fmt.Errorf("CSV deserialization requires slice of structs, got slice of %s", elemType.Kind())
	}

	// Build field map: CSV column name -> struct field index
	fieldMap := make(map[string]int)
	for i := 0; i < elemType.NumField(); i++ {
		field := elemType.Field(i)
		if filter.IsFieldIncluded(field, groups) {
			fieldMap[field.Name] = i
		}
	}

	// Process data rows
	newSlice := reflect.MakeSlice(sliceValue.Type(), 0, len(rows)-1)
	for i := 1; i < len(rows); i++ {
		row := rows[i]
		elem := reflect.New(elemType).Elem()

		for j, cellValue := range row {
			if j >= len(header) {
				break
			}
			fieldName := header[j]
			fieldIdx, ok := fieldMap[fieldName]
			if !ok {
				continue
			}

			field := elem.Field(fieldIdx)
			if !field.CanSet() {
				continue
			}

			// Set field value based on type
			if err := setFieldFromString(field, cellValue); err != nil {
				return fmt.Errorf("error setting field %s: %w", fieldName, err)
			}
		}

		if isPtr {
			newSlice = reflect.Append(newSlice, elem.Addr())
		} else {
			newSlice = reflect.Append(newSlice, elem)
		}
	}

	sliceValue.Set(newSlice)
	return nil
}

func setFieldFromString(field reflect.Value, value string) error {
	switch field.Kind() {
	case reflect.String:
		field.SetString(value)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if value == "" {
			value = "0"
		}
		intVal := int64(0)
		_, err := fmt.Sscanf(value, "%d", &intVal)
		if err != nil {
			return err
		}
		field.SetInt(intVal)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if value == "" {
			value = "0"
		}
		uintVal := uint64(0)
		_, err := fmt.Sscanf(value, "%d", &uintVal)
		if err != nil {
			return err
		}
		field.SetUint(uintVal)
	case reflect.Float32, reflect.Float64:
		if value == "" {
			value = "0"
		}
		floatVal := float64(0)
		_, err := fmt.Sscanf(value, "%f", &floatVal)
		if err != nil {
			return err
		}
		field.SetFloat(floatVal)
	case reflect.Bool:
		boolVal := value == "true" || value == "1"
		field.SetBool(boolVal)
	default:
		return fmt.Errorf("unsupported field type: %s", field.Kind())
	}
	return nil
}
