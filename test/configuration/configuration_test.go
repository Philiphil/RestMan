package configuration_test

import (
	"testing"

	. "github.com/philiphil/restman/configuration"
)

// test Configuration.GetValueAsString

func TestConfiguration_GetValueAsString(t *testing.T) {
	config := Configuration{}
	config.Values = []string{"value1", "value2", "value3"}

	value, err := config.GetValueAsString()
	if err != nil {
		t.Error(err)
	}
	if value != "value1" {
		t.Error("unexpected value:", value)
	}

	//now we empty it
	config.Values = config.Values[:0]
	value, err = config.GetValueAsString()
	if err == nil {
		t.Error(err)
	}
	if value != "" {
		t.Error("expected empty value but got:", value)
	}
}

func TestConfiguration_GetValueAsBool(t *testing.T) {
	config := Configuration{}
	config.Values = []string{"true", "false"}

	value, err := config.GetValueAsBool()
	if err != nil {
		t.Error(err)
	}
	if value != true {
		t.Error("unexpected value:", value)
	}
	//now use 1 0
	config.Values = []string{"1", "0"}
	value, err = config.GetValueAsBool()
	if err != nil {
		t.Error(err)
	}
	if value != true {
		t.Error("unexpected value:", value)
	}
	//now use bad value
	config.Values = []string{"yes", "no"}
	value, err = config.GetValueAsBool()
	if err == nil {
		t.Error("expected error but got none")
	}
	if value != false {
		t.Error("expected false value but got:", value)
	}
	//now empty
	config.Values = config.Values[:0]
	value, err = config.GetValueAsBool()
	if err == nil {
		t.Error("expected error but got none")
	}
	if value != false {
		t.Error("expected false value but got:", value)
	}
}

func TestConfiguration_GetValusAsInt(t *testing.T) {
	config := Configuration{}
	config.Values = []string{"42", "100"}

	value, err := config.GetValueAsInt()
	if err != nil {
		t.Error(err)
	}
	if value != 42 {
		t.Error("unexpected value:", value)
	}
	//now use bad value
	config.Values = []string{"forty-two", "one hundred"}
	value, err = config.GetValueAsInt()
	if err == nil {
		t.Error("expected error but got none")
	}
	if value != 0 {
		t.Error("expected 0 value but got:", value)
	}
	//now empty
	config.Values = config.Values[:0]
	value, err = config.GetValueAsInt()
	if err == nil {
		t.Error("expected error but got none")
	}
	if value != 0 {
		t.Error("expected 0 value but got:", value)
	}
}
