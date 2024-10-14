package serializer_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/philiphil/restman/format"
	. "github.com/philiphil/restman/serializer"
)

type Test struct {
	Test0 int `group:"test"`
	Test1 int `group:"testo"`
	Test2 int `group:"test"`
	Test3 int `group:"testo,test"`
	Test4 int
	test5 int
	Test6 int `group:"test"`
}

type Recursive struct {
	Test1 Hidden `group:"test"`
	Test2 Hidden
}
type Hidden struct {
	Test0 int `group:"test"`
	Test1 int
}

type Ptr struct {
	Test0 int     `group:"test"`
	Test1 *int    `group:"test"`
	Test2 *Hidden `group:"test"`
	Test3 *int
	Test4 *Hidden
}

type Test2 struct {
	Test
}

var test = Test{
	9, -8, 7, 6, -5, -4, 3,
}
var testDeserializedResult = Test{
	9, 0, 7, 6, 0, 0, 3,
}

// basic struct
func TestSerializer_Deserialize(t *testing.T) {
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
	if o != testDeserializedResult || err != nil {
		t.Error("!")
	}

}

// nested struct
func TestSerializer_Deserialize2(t *testing.T) {
	test2 := Recursive{
		Hidden{1, 2},
		Hidden{3, 4},
	}
	expected2 := Recursive{
		Hidden{1, 0},
		Hidden{0, 0},
	}

	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		t.Error(err)
	}
	o := Recursive{}
	err = s.Deserialize(serialized, &o)
	if o != expected2 || err != nil {
		fmt.Println(o)
		fmt.Println(expected2)
		t.Error("!")
	}

}

// slice
func TestSerializer_Deserialize3(t *testing.T) {
	test2 := []Recursive{
		{
			Hidden{1, 2},
			Hidden{3, 4},
		},
	}
	expected2 := []Recursive{
		{
			Hidden{1, 0},
			Hidden{0, 0},
		},
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		t.Error(err)
	}
	o := []Recursive{}
	err = s.Deserialize(serialized, &o)
	if o[0] != expected2[0] || err != nil {
		fmt.Println(o)
		fmt.Println(expected2)
		t.Error("!")
	}

}

// ptr to struct
func TestSerializer_Deserialize4(t *testing.T) {
	test1 := new(Hidden)
	expected1 := new(Hidden)
	test1.Test0 = 1
	test1.Test1 = 1
	expected1.Test0 = 1
	expected1.Test1 = 0

	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test1, "test")
	if err != nil {
		t.Error(err)
	}
	o := new(Hidden)
	err = s.Deserialize(serialized, &o)
	if *o != *expected1 || err != nil {
		t.Error("!")
	}

}

// slice of ptr
func TestSerializer_Deserialize5(t *testing.T) {
	test2 := []*Recursive{
		{
			Hidden{1, 2},
			Hidden{3, 4},
		},
	}
	expected2 := []*Recursive{
		{
			Hidden{1, 0},
			Hidden{0, 0},
		},
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		t.Error(err)
	}
	o := []*Recursive{}
	err = s.Deserialize(serialized, &o)
	if *o[0] != *expected2[0] || err != nil {
		fmt.Println(*o[0])
		fmt.Println(*expected2[0])
		t.Error("!")
	}

}

// ptr slice
func TestSerializer_Deserialize6(t *testing.T) {
	test2 := new([]Recursive)
	*test2 = []Recursive{
		{
			Hidden{1, 2},
			Hidden{3, 4},
		},
	}
	expected2 := []Recursive{
		{
			Hidden{1, 0},
			Hidden{0, 0},
		},
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		t.Error(err)
	}
	o := new([]Recursive)
	err = s.Deserialize(serialized, o)
	if (*o)[0] != expected2[0] || err != nil {
		fmt.Println(*o)
		fmt.Println(expected2)
		t.Error("!")
	}
}

// struct with ptr
func TestSerializer_Deserialize7(t *testing.T) {
	test2 := new([]Recursive)
	*test2 = []Recursive{
		{
			Hidden{1, 2},
			Hidden{3, 4},
		},
	}
	expected2 := []Recursive{
		{
			Hidden{1, 0},
			Hidden{0, 0},
		},
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		t.Error(err)
	}
	o := new([]Recursive)
	err = s.Deserialize(serialized, o)
	if (*o)[0] != expected2[0] || err != nil {
		fmt.Println(*o)
		fmt.Println(expected2)
		t.Error("!")
	}
}

// struct w/ ptr & ptr to nested
func TestSerializer_Deserialize8(t *testing.T) {
	intValue := 42

	test := Ptr{
		Test0: 1,
		Test1: &intValue,
		Test2: &Hidden{2, 3},
		Test3: &intValue,
		Test4: &Hidden{2, 3},
	}

	expected := Ptr{
		Test0: 1,
		Test1: &intValue,
		Test2: &Hidden{2, 0},
		Test3: nil,
		Test4: nil,
	}

	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}

	var o Ptr
	err = s.Deserialize(serialized, &o)
	if !reflect.DeepEqual(o, expected) || err != nil {
		fmt.Println(o)
		fmt.Println(expected)
		t.Error("Test failed!")
	}
}

// map[any] to map[typed]
func TestSerializer_Deserialize9(t *testing.T) {
	test1 := make(map[string]any)
	test1["test"] = test
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test1, "test")
	if err != nil {
		t.Error(err)
	}
	o := make(map[string]Test)
	err = s.Deserialize(serialized, &o)
	if o["test"] != testDeserializedResult || err != nil {
		fmt.Println(o["test"])
		fmt.Println(testDeserializedResult)
		t.Error("!")
	}
}

// anonymous
func TestSerializer_Deserialize11(t *testing.T) {
	s := NewSerializer(format.JSON)
	tt := Test2{test}
	rr := Test2{testDeserializedResult}
	serialized, err := s.Serialize(tt, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test2{}
	err = s.Deserialize(serialized, &o)
	if o != rr || err != nil {
		fmt.Println(serialized)
		fmt.Println(rr)
		fmt.Println(o)
		t.Error("!")
	}

}

// anonymous w/ptr
func TestSerializer_Deserialize12(t *testing.T) {
	s := NewSerializer(format.JSON)
	tt := Test2{test}
	rr := Test2{testDeserializedResult}
	serialized, err := s.Serialize(&tt, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test2{}
	err = s.Deserialize(serialized, &o)
	if o != rr || err != nil {
		fmt.Println(serialized)
		fmt.Println(rr)
		fmt.Println(o)
		t.Error("!")
	}

}

func TestSerializer_MergeObjects(t *testing.T) {
	target := Test{
		11, 11, 11, 11, 11, 11, 11,
	}
	result := Test{
		9, 11, 7, 6, 11, 11, 3,
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
	if err != nil {
		t.Error(err)
	}
	err = s.MergeObjects(&target, &o)
	if err != nil {
		t.Error(err)
	}
	if target != result {
		fmt.Println(target)
		fmt.Println(result)
		t.Error("!")
	}
	//object should be a pointer
	err = s.MergeObjects(target, &o)
	if err == nil {
		t.Error(err)
	}

	err = s.MergeObjects(&target, o)
	if err == nil {
		t.Error(err)
	}

}

func TestSerializer_DeserializeAndMerge(t *testing.T) {
	target := Test{
		11, 11, 11, 11, 11, 11, 11,
	}
	result := Test{
		9, 11, 7, 6, 11, 11, 3,
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}
	err = s.DeserializeAndMerge(serialized, &target)
	if err != nil {
		t.Error(err)
	}
	if target != result {
		fmt.Println(target)
		fmt.Println(result)
		t.Error("!")
	}

}

// target object should be a pointer
func TestSerializer_DeserializeStruct(t *testing.T) {
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, o)
	if err == nil {
		t.Error("!")
	}

}

// jsonLd coverage
func TestSerializer_SerializeJsonLD(t *testing.T) {
	s := NewSerializer(format.JSONLD)
	serialized, err := s.Serialize(test, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
	if o != testDeserializedResult || err != nil {
		t.Error("!")
	}

}

// Undefined format coverage
func TestSerializer_SerializeUndefined(t *testing.T) {
	s := NewSerializer(format.Undefined)
	_, err := s.Serialize(test, "test")
	if err == nil {
		t.Error("should not work")
	}

}

// Unknown format coverage
func TestSerializer_SerializeUnknown(t *testing.T) {
	s := NewSerializer(format.Unknown)
	_, err := s.Serialize(test, "test")
	if err == nil {
		t.Error(err)
	}
}
