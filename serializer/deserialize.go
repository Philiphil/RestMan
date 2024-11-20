package serializer

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/serializer/filter"
)

// Deserializer encapsulates the deserialization logic
func (s *Serializer) Deserialize(data string, obj any) error {
	if !isPointer(obj) {
		return fmt.Errorf("object must be pointer")
	}
	switch s.Format {
	case format.JSONLD:
		return json.Unmarshal([]byte(data), obj)
	case format.JSON:
		return json.Unmarshal([]byte(data), obj)
	//case format.XML:
	//	return xml.Unmarshal([]byte(data), obj)
	//case format.CSV:
	//	return s.deserializeCSV(data, obj)
	default:
		return fmt.Errorf("unsupported format: %s", s.Format)
	}
}

// MergeObjects merges two objects together
// Both target and source must be pointers
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

func (s *Serializer) deserializeCSV(data string, obj any) error {
	value := reflect.ValueOf(obj)
	if value.Kind() != reflect.Ptr || value.IsNil() {
		return fmt.Errorf("invalid object type for CSV deserialization")
	}

	reader := csv.NewReader(strings.NewReader(data))
	rows, err := reader.ReadAll()
	if err != nil {
		return err
	}

	elemType := value.Elem().Type()
	elem := reflect.New(elemType).Elem()

	for _, row := range rows {
		for i, fieldValue := range row {
			field := elem.Field(i)
			if field.IsValid() && field.CanSet() {
				field.SetString(fieldValue)
			}
		}
	}

	value.Elem().Set(elem)
	return nil
}
