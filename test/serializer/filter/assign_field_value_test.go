package filter

import (
	"reflect"
	"testing"

	. "github.com/philiphil/restman/serializer/filter"
)

// simple struct for nested copy test
type inner struct {
	A int
	B string
}

type outer struct {
	I inner
}

// TestAssignFieldValueIncompatible ensures an error is returned when types are not convertible.
func TestAssignFieldValueIncompatible(t *testing.T) {
	var target int
	// source is a string; cannot convert to int via Convert unless string numeric; our function checks ConvertibleTo.
	src := reflect.ValueOf("abc")
	field := reflect.StructField{Type: reflect.TypeOf(target)}
	destVal := reflect.New(reflect.TypeOf(target)).Elem()
	if err := AssignFieldValue(field, destVal, src); err == nil {
		t.Fatalf("expected error for non convertible types")
	}
}

// TestAssignFieldValueNestedStruct copies nested struct values.
func TestAssignFieldValueNestedStruct(t *testing.T) {
	src := outer{I: inner{A: 42, B: "x"}}
	var dst outer

	srcVal := reflect.ValueOf(src)
	dstVal := reflect.ValueOf(&dst).Elem()

	// iterate fields manually to exercise recursion
	for i := 0; i < dstVal.NumField(); i++ {
		field := dstVal.Type().Field(i)
		if err := AssignFieldValue(field, dstVal.Field(i), srcVal.Field(i)); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}
	if dst.I.A != 42 || dst.I.B != "x" {
		t.Fatalf("expected nested struct copy, got %+v", dst.I)
	}
}
