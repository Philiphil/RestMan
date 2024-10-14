package router_test

import (
	"testing"

	. "github.com/philiphil/restman/router"
)

func TestConvertToSnakeCase(t *testing.T) {
	if ConvertToSnakeCase("") != "" {
		t.Error("!")
	}
	if ConvertToSnakeCase("Test") != "test" {
		t.Error("!")
	}
	if ConvertToSnakeCase("TestTest") != "test_test" {
		t.Error("!")
	}
	if ConvertToSnakeCase("TestTestTest") != "test_test_test" {
		t.Error("!")
	}
	if ConvertToSnakeCase("TestTestTestTest") != "test_test_test_test" {
		t.Error("!")
	}
}
