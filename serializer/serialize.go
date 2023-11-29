package serializer

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/philiphil/apiman/serializer/filter"
	"github.com/philiphil/apiman/serializer/format"
	"reflect"
	"strings"
)

func (s *Serializer) Serialize(obj any, groups ...string) (string, error) {
	switch s.Format {
	case format.JSON:
		return s.serializeJSON(obj, groups...)
	case format.XML:
		return s.serializeXML(obj, groups...)
	case format.CSV:
		return s.serializeCSV(obj, groups...)
	default:
		return "", fmt.Errorf("Unsupported format: %s", s.Format)
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
	data := filter.FilterByGroups(obj, groups...)
	xmlBytes, err := xml.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(xmlBytes), nil
}

func (s *Serializer) serializeCSV(obj any, groups ...string) (string, error) {
	data := filter.FilterByGroups(obj, groups...)

	value := reflect.ValueOf(data)
	if value.Kind() != reflect.Slice {
		return "", fmt.Errorf("Invalid object type for CSV serialization")
	}

	rows := make([][]string, 0)
	header := make([]string, 0)

	for i := 0; i < value.Len(); i++ {
		row := make([]string, 0)
		elem := value.Index(i).Interface()
		elemValue := reflect.ValueOf(elem)

		// Handle header row
		if i == 0 {
			for j := 0; j < elemValue.NumField(); j++ {
				field := elemValue.Type().Field(j)
				if filter.IsFieldIncluded(field, groups) {
					header = append(header, field.Name)
				}
			}
			rows = append(rows, header)
		}

		for j := 0; j < elemValue.NumField(); j++ {
			field := elemValue.Field(j)
			if filter.IsFieldIncluded(elemValue.Type().Field(j), groups) {
				row = append(row, fmt.Sprintf("%v", field.Interface()))
			}
		}

		rows = append(rows, row)
	}

	csvBytes, err := writeCSVToString(rows)
	if err != nil {
		return "", err
	}

	return string(csvBytes), nil
}

func writeCSVToString(rows [][]string) ([]byte, error) {
	sb := strings.Builder{}
	writer := csv.NewWriter(&sb)
	err := writer.WriteAll(rows)
	if err != nil {
		return nil, err
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return []byte(sb.String()), nil
}
