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

type xmlFilterWrapper struct {
	Data any
	groups []string
}

func (w xmlFilterWrapper) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	return marshalXMLFiltered(e, start, reflect.ValueOf(w.Data), w.groups)
}

func marshalXMLFiltered(e *xml.Encoder, start xml.StartElement, value reflect.Value, groups []string) error {
	value = filter.DereferenceValueIfPointer(value)

	if value.Kind() == reflect.Slice {
		// Encode opening tag for slice container
		if err := e.EncodeToken(start); err != nil {
			return err
		}

		// Encode each element
		for i := 0; i < value.Len(); i++ {
			elem := value.Index(i)
			elemStart := xml.StartElement{Name: xml.Name{Local: "item"}}
			if err := marshalXMLFiltered(e, elemStart, elem, groups); err != nil {
				return err
			}
		}

		// Encode closing tag
		return e.EncodeToken(start.End())
	}

	if value.Kind() != reflect.Struct {
		// For primitives, just encode directly
		if err := e.EncodeElement(value.Interface(), start); err != nil {
			return err
		}
		return nil
	}

	// For structs, filter fields by group
	if err := e.EncodeToken(start); err != nil {
		return err
	}

	typ := value.Type()
	for i := 0; i < value.NumField(); i++ {
		field := typ.Field(i)
		fieldValue := value.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Check if field should be included based on groups
		if !filter.IsFieldIncluded(field, groups) {
			continue
		}

		// Get XML tag name
		xmlTag := field.Tag.Get("xml")
		fieldName := field.Name
		if xmlTag != "" && xmlTag != "-" {
			parts := strings.Split(xmlTag, ",")
			if parts[0] != "" {
				fieldName = parts[0]
			}
		}

		fieldStart := xml.StartElement{Name: xml.Name{Local: fieldName}}

		if filter.IsStruct(field.Type) {
			if err := marshalXMLFiltered(e, fieldStart, fieldValue, groups); err != nil {
				return err
			}
		} else {
			if err := e.EncodeElement(fieldValue.Interface(), fieldStart); err != nil {
				return err
			}
		}
	}

	return e.EncodeToken(start.End())
}

func (s *Serializer) Serialize(obj any, groups ...string) (string, error) {
	switch s.Format {
	case format.JSON:
		return s.serializeJSON(obj, groups...)
	case format.JSONLD:
		return s.serializeJSON(obj, groups...)
	case format.XML:
		return s.serializeXML(obj, groups...)
	case format.CSV:
		return s.serializeCSV(obj, groups...)
	default:
		return "", fmt.Errorf("unsupported format: %s", s.Format)
	}
}

func (s *Serializer) serializeJSON(obj any, groups ...string) (string, error) {
	data := filter.FilterByGroups(obj, groups...)
	//data = renameFieldsToLower(data)
	jsonBytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func (s *Serializer) serializeXML(obj any, groups ...string) (string, error) {
	// Use custom XML marshaler that filters by groups
	var sb strings.Builder
	sb.WriteString(xml.Header)

	encoder := xml.NewEncoder(&sb)
	encoder.Indent("", "  ")

	value := reflect.ValueOf(obj)
	typ := value.Type()

	// Determine root element name
	rootName := "root"
	if typ.Kind() == reflect.Struct {
		rootName = typ.Name()
		if rootName == "" {
			rootName = "root"
		}
	} else if typ.Kind() == reflect.Slice {
		rootName = "items"
	}

	start := xml.StartElement{Name: xml.Name{Local: rootName}}
	if err := marshalXMLFiltered(encoder, start, value, groups); err != nil {
		return "", err
	}

	if err := encoder.Flush(); err != nil {
		return "", err
	}

	return sb.String(), nil
}

func (s *Serializer) serializeCSV(obj any, groups ...string) (string, error) {
	data := filter.FilterByGroups(obj, groups...)

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		return "", fmt.Errorf("invalid object type for CSV serialization: expected slice, got %s", value.Kind())
	}

	if value.Len() == 0 {
		return "", nil
	}

	// Build header and field indices once
	firstElem := value.Index(0)
	elemValue := reflect.ValueOf(firstElem.Interface())
	if elemValue.Kind() == reflect.Ptr {
		elemValue = elemValue.Elem()
	}

	if elemValue.Kind() != reflect.Struct {
		return "", fmt.Errorf("CSV serialization requires slice of structs, got slice of %s", elemValue.Kind())
	}

	// Cache field indices and build header
	var fieldIndices []int
	var header []string
	for j := 0; j < elemValue.NumField(); j++ {
		field := elemValue.Type().Field(j)
		if filter.IsFieldIncluded(field, groups) {
			fieldIndices = append(fieldIndices, j)
			header = append(header, field.Name)
		}
	}

	if len(header) == 0 {
		return "", fmt.Errorf("no fields to serialize in CSV")
	}

	// Pre-allocate rows slice
	rows := make([][]string, 0, value.Len()+1)
	rows = append(rows, header)

	// Process each element
	for i := 0; i < value.Len(); i++ {
		elem := value.Index(i)
		elemValue := reflect.ValueOf(elem.Interface())
		if elemValue.Kind() == reflect.Ptr {
			elemValue = elemValue.Elem()
		}

		row := make([]string, 0, len(fieldIndices))
		for _, idx := range fieldIndices {
			field := elemValue.Field(idx)
			row = append(row, fmt.Sprintf("%v", field.Interface()))
		}
		rows = append(rows, row)
	}

	return writeCSVToString(rows)
}

func writeCSVToString(rows [][]string) (string, error) {
	sb := strings.Builder{}
	writer := csv.NewWriter(&sb)
	if err := writer.WriteAll(rows); err != nil {
		return "", err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return "", err
	}
	return sb.String(), nil
}
