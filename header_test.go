package apiman_test

import (
	"testing"

	_ "github.com/philiphil/apiman"
	_ "github.com/philiphil/apiman/format"
	"github.com/philiphil/apiman/router"
)

func TestParseAcceptHeader(t *testing.T) {
	_, err := router.ParseAcceptHeader("text/plain; q=0.5, text/html,text/x-dvi; q=0.8, text/x-c")
	if err == nil {
		//nothing acceptable sent
		t.Error(err)
	}
	_, err = router.ParseAcceptHeader("*/*")
	if err != nil {
		//*/* should be acceptable
		t.Error(err)
	}
	_, err = router.ParseAcceptHeader("application/json")
	if err != nil {
		//application/json should be acceptable
		t.Error(err)
	}

}
