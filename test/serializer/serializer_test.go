package serializer_test

import (
	_ "github.com/philiphil/restman/serializer"
)

// needs fix

/*
func TestSerializer_SerializeXML(t *testing.T) {
	s := NewSerializer(format.XML)
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

// needs fix
func TestSerializer_SerializeCSV(t *testing.T) {
	s := NewSerializer(format.CSV)
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

// slice
func TestSerializer_Deserialize13(t *testing.T) {
	s := NewSerializer(format.JSON)
	serialized, err := s.Serialize([]Test{test}, "test")
	if err != nil {
		t.Error(err)
	}
	o := Test{}
	err = s.Deserialize(serialized, &o)
	if o != testDeserializedResult || err != nil {
		t.Error("!")
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
*/
