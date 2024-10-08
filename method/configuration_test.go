package method_test

import (
	"testing"

	. "github.com/philiphil/restman/method"
)

func TestDefaultApiMethods(t *testing.T) {
	test := DefaultApiMethods()
	if len(test) != 8 {
		t.Error("DefaultApiMethods should return 8 methods")
	}

}
