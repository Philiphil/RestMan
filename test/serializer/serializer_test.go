package serializer_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/philiphil/restman/format"
	. "github.com/philiphil/restman/serializer"
)

// XML serialization - basic struct
func TestSerializer_SerializeXML(t *testing.T) {
	s := NewSerializer(format.XML)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}

	// Check XML header is present
	if !strings.HasPrefix(serialized, "<?xml") {
		t.Error("XML should start with header")
	}

	o := Test{}
	err = s.Deserialize(serialized, &o)
	if o != testDeserializedResult || err != nil {
		fmt.Println("Got:", o)
		fmt.Println("Expected:", testDeserializedResult)
		t.Error("XML deserialization failed")
	}
}

// XML serialization - slice of structs
func TestSerializer_SerializeXMLSlice(t *testing.T) {
	s := NewSerializer(format.XML)
	testSlice := []Test{test, test}
	serialized, err := s.Serialize(testSlice, "test")
	if err != nil {
		t.Error(err)
	}

	o := []Test{}
	err = s.Deserialize(serialized, &o)
	if err != nil {
		t.Error(err)
	}
	if len(o) != 2 || o[0] != testDeserializedResult {
		fmt.Println("Got:", o)
		t.Error("XML slice deserialization failed")
	}
}

// CSV serialization - slice of structs
func TestSerializer_SerializeCSV(t *testing.T) {
	s := NewSerializer(format.CSV)
	testSlice := []Test{test, test}
	serialized, err := s.Serialize(testSlice, "test")
	if err != nil {
		t.Error(err)
	}

	// Check CSV has header
	lines := strings.Split(strings.TrimSpace(serialized), "\n")
	if len(lines) < 2 {
		t.Error("CSV should have at least header and one data row")
	}

	o := []Test{}
	err = s.Deserialize(serialized, &o)
	if err != nil {
		t.Error(err)
	}
	if len(o) != 2 || o[0] != testDeserializedResult {
		fmt.Println("Got:", o)
		fmt.Println("Expected:", testDeserializedResult)
		t.Error("CSV deserialization failed")
	}
}

// CSV serialization - empty slice
func TestSerializer_SerializeCSVEmpty(t *testing.T) {
	s := NewSerializer(format.CSV)
	testSlice := []Test{}
	serialized, err := s.Serialize(testSlice, "test")
	if err != nil {
		t.Error(err)
	}
	if serialized != "" {
		t.Error("Empty slice should produce empty CSV")
	}
}

// CSV serialization - slice with different types
func TestSerializer_SerializeCSVPrimitives(t *testing.T) {
	type CSVTest struct {
		Name  string  `groups:"test"`
		Age   int     `groups:"test"`
		Score float64 `groups:"test"`
		Pass  bool    `groups:"test"`
	}

	s := NewSerializer(format.CSV)
	testData := []CSVTest{
		{"Alice", 25, 95.5, true},
		{"Bob", 30, 87.3, false},
	}

	serialized, err := s.Serialize(testData, "test")
	if err != nil {
		t.Error(err)
	}

	o := []CSVTest{}
	err = s.Deserialize(serialized, &o)
	if err != nil {
		t.Error(err)
	}

	if len(o) != 2 {
		t.Errorf("Expected 2 items, got %d", len(o))
	}
	if o[0].Name != "Alice" || o[0].Age != 25 || o[0].Score != 95.5 || o[0].Pass != true {
		fmt.Println("Got:", o[0])
		t.Error("CSV primitive types deserialization failed")
	}
}

// CSV error handling - not a slice
func TestSerializer_SerializeCSVNotSlice(t *testing.T) {
	s := NewSerializer(format.CSV)
	_, err := s.Serialize(test, "test")
	if err == nil {
		t.Error("Should error when serializing non-slice to CSV")
	}
}

// CSV error handling - slice of non-structs
func TestSerializer_SerializeCSVNonStruct(t *testing.T) {
	s := NewSerializer(format.CSV)
	testSlice := []int{1, 2, 3}
	_, err := s.Serialize(testSlice, "test")
	if err == nil {
		t.Error("Should error when serializing slice of primitives to CSV")
	}
}

// map[typed] to map[typed]
func TestSerializer_Deserialize10(t *testing.T) {
	test1 := make(map[string]Test)
	test1["test"] = test
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test1, "test")
	if err != nil {
		t.Error(err)
	}
	o := make(map[string]Test)
	err = s.Deserialize(serialized, &o)
	if err != nil {
		t.Error(err)
	}
	if o["test"] != testDeserializedResult {
		fmt.Println("o[test]:", o["test"])
		fmt.Println("Expected:", testDeserializedResult)
		t.Error("Deserialization result does not match expected value")
	}
}
