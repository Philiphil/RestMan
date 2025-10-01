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

// slice
func TestSerializer_Deserialize13(t *testing.T) {
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize([]Test{test}, "test")
	if err != nil {
		t.Error(err)
	}
	o := []Test{}
	err = s.Deserialize(serialized, &o)
	if o[0] != testDeserializedResult || err != nil {
		fmt.Println("o:", o)
		fmt.Println("Expected:", testDeserializedResult)
		t.Error("!")
	}
}

// try to deserialize primitive type
func TestSerializer_Deserialize14(t *testing.T) {
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(1, "test")
	if err != nil {
		t.Error(err)
	}
	var o int
	err = s.Deserialize(serialized, &o)
	if o != 1 || err != nil {
		t.Error("!")
	}
	//now string
	serialized, err = s.Serialize("test", "test")
	if err != nil {
		t.Error(err)
	}
	var o2 string
	err = s.Deserialize(serialized, &o2)
	if o2 != "test" || err != nil {
		t.Error("!")
	}
	//now bool
	serialized, err = s.Serialize(true, "test")
	if err != nil {
		t.Error(err)
	}
	var o3 bool
	err = s.Deserialize(serialized, &o3)
	if o3 != true || err != nil {
		t.Error("!")
	}
	//now float
	serialized, err = s.Serialize(1.1, "test")
	if err != nil {
		t.Error(err)
	}
	var o4 float64
	err = s.Deserialize(serialized, &o4)
	if o4 != 1.1 || err != nil {
		t.Error("!")
	}
	// now int64
	serialized, err = s.Serialize(int64(9223372036854775806), "test")
	if err != nil {
		t.Error(err)
	}
	var o5 int64
	err = s.Deserialize(serialized, &o5)
	if o5 != int64(9223372036854775806) || err != nil {
		t.Error("!")
	}
	// now Rune
	serialized, err = s.Serialize('a', "test")
	if err != nil {
		t.Error(err)
	}
	var o6 rune
	err = s.Deserialize(serialized, &o6)
	if o6 != 'a' || err != nil {
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

// test performance of deserialization, use a slice of 100 PTR structs,
func BenchmarkSerializer_Deserialize(b *testing.B) {
	test2 := []*Recursive{}
	for i := 0; i < 100; i++ {
		test2 = append(test2, &Recursive{
			Hidden{1, 2},
			Hidden{3, 4},
		})
	}
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize(test2, "test")
	if err != nil {
		b.Error(err)
	}
	o := []*Recursive{}
	for i := 0; i < b.N; i++ {
		err = s.Deserialize(serialized, &o)
		if err != nil {
			b.Error(err)
		}
	}
}

// test performance of serialization, use a slice of 100 Test2 structs
func BenchmarkSerializer_Serialize(b *testing.B) {
	test2 := []Test2{}
	for i := 0; i < 100; i++ {
		test2 = append(test2, Test2{test})
	}
	s := NewSerializer(format.JSON)
	for i := 0; i < b.N; i++ {
		_, err := s.Serialize(test2, "test")
		if err != nil {
			b.Error(err)
		}
	}
}

// todo we need better test
type FilterStruct struct {
	Test0 int `group:"test"`
	Test1 int `group:"testo"`
	Test2 int
	Test3 int `group:"testo,test"`
}

var testedFilterStruct = FilterStruct{
	9, -8, 7, 6,
}

var testedFilterStructDeserializedResult = FilterStruct{
	9, 0, 7, 6,
}

type PrimitiveStruct struct {
	Test0 int     `group:"test"`
	Test1 string  `group:"test"`
	Test2 bool    `group:"test"`
	Test3 float64 `group:"test"`
	Test4 int64   `group:"test"`
	Test5 rune    `group:"test"`
}

type NestedStruct struct {
	Test0 FilterStruct `group:"test"`
	Test1 FilterStruct
}

type PtrStruct struct {
	Test0 *FilterStruct `group:"test"`
	Test1 *FilterStruct
}

type AnonymousStruct struct {
	FilterStruct
}

type SliceStruct struct {
	Test0 []FilterStruct `group:"test"`
	Test1 []FilterStruct
}

type MapStruct struct {
	Test0 map[string]FilterStruct `group:"test"`
	Test1 map[string]FilterStruct
}

type InterfaceStruct struct {
	Test0 any `group:"test"`
	Test1 any
}
