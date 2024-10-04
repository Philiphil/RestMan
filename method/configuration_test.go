package method_test

import (
	"testing"

	. "github.com/philiphil/apiman/method"
)

func TestDefaultApiMethods(t *testing.T) {
	test := DefaultApiMethods()
	if len(test) != 10 {
		t.Error("DefaultApiMethods should return 10 methods")
	}

}
