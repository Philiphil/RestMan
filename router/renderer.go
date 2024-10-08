package router

import (
	"net/http"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/serializer"
)

type SerializerRenderer struct {
	Data   any
	Format format.Format
	Groups []string
}

var (
	jsonContentType   = []string{"application/json; charset=utf-8"}
	jsonldContentType = []string{"application/ld-json; charset=utf-8"}
	//xmlContentType    = []string{"application/xml; charset=utf-8"}
	//csvContentType    = []string{"text/csv"}
)

// Render
func (r SerializerRenderer) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	s := serializer.NewSerializer(r.Format)
	str, err := s.Serialize(r.Data, r.Groups...)
	if err != nil {
		panic(err)
	}
	_, err = w.Write([]byte(str))
	return err
}

// WriteContentType (JSON) writes JSON ContentType.
func (r SerializerRenderer) WriteContentType(w http.ResponseWriter) {
	switch r.Format {
	case format.JSON:
		writeContentType(w, jsonContentType)
	case format.JSONLD:
		writeContentType(w, jsonldContentType)
	case format.XML:
		//	writeContentType(w, xmlContentType)
	case format.CSV:
		//	writeContentType(w, csvContentType)
	}
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
