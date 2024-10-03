package router_test

import (
	"testing"

	"github.com/philiphil/apiman/format"
	. "github.com/philiphil/apiman/router"
)

func TestParseAcceptHeader(t *testing.T) {
	if _, err := ParseAcceptHeader("text/plain; q=0.5, text/html,text/x-dvi; q=0.8, text/x-c"); err == nil {
		t.Error(err)
	}
	if format_, err := ParseAcceptHeader("*/*"); err != nil || format_ != format.JSON {
		t.Error(err)
	}
	if format_, err := ParseAcceptHeader("text/plain; q=0.5, text/html,text/x-dvi; q=0.8, application/json"); err != nil || format_ != format.JSON {
		t.Error(err)
	}

	if format_, err := ParseAcceptHeader("text/plain; q=0.5, text/html,application/x-json; q=0.8, text/x-c"); err != nil || format_ != format.JSON {
		t.Error(err)
	}

	if format_, err := ParseAcceptHeader(""); err != nil || format_ != format.JSON {
		t.Error(err)
	}

}

func TestParseTypeFromString(t *testing.T) {
	if ParseTypeFromString("") != format.Undefined {
		t.Error("Expected Undefined")
	}
	if ParseTypeFromString("application/json") != format.JSON {
		t.Error("Expected JSON")
	}
	if ParseTypeFromString("application/xml") != format.XML {
		t.Error("Expected XML")
	}
	if ParseTypeFromString("application/ld+json") != format.JSONLD {
		t.Error("Expected JSONLD")
	}
	if ParseTypeFromString("application/csv") != format.CSV {
		t.Error("Expected csv")
	}
	if ParseTypeFromString("this one is a new one") != format.Unknown {
		t.Error("Expected Unknown")
	}

}
