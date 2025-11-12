package router_test

import (
	"testing"

	"github.com/philiphil/restman/format"
	. "github.com/philiphil/restman/router"
)

// TestParseAcceptHeaderWeighted ensures q weights and tie-break alphabetical ordering work.
func TestParseAcceptHeaderWeighted(t *testing.T) {
	header := "text/csv;q=0.5, application/xml;q=0.9, application/json;q=0.9"
	f, err := ParseAcceptHeader(header)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// xml and json have same q=0.9; json < xml lexicographically so JSON should win
	if f != format.JSON {
		t.Fatalf("expected JSON, got %v", f)
	}
}

// TestParseAcceptHeaderWildcard falls back to JSON when */* present.
func TestParseAcceptHeaderWildcard(t *testing.T) {
	f, err := ParseAcceptHeader("*/*")
	if err != nil || f != format.JSON {
		t.Fatalf("expected JSON for */*, got %v (err=%v)", f, err)
	}
}
