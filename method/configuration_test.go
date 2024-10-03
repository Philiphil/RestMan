package method_test

import (
	"testing"

	. "github.com/philiphil/apiman/method"
)

func TestDefaultApiMethods(t *testing.T) {
	test := DefaultApiMethods()
	if len(test) != 9 {
		t.Error("DefaultApiMethods should return 9 methods")
	}

}
