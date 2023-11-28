package Gin

import (
	"ApiMan/Serializer/Format"
	"net/http"
)

type SerializerRenderer struct {
	Data   any
	Format Format.Format
	Groups []string
}

var (
	jsonContentType = []string{"application/json; charset=utf-8"}
	xmlContentType  = []string{"application/xml; charset=utf-8"}
	csvContentType  = []string{"text/csv"}
)

// Render
func (r SerializerRenderer) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	s := serializer.NewSerializer(r.Format)
	str, err := s.Serialize(r.Data, r.Groups...)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(str))
	return err
}

// WriteContentType (JSON) writes JSON ContentType.
func (r SerializerRenderer) WriteContentType(w http.ResponseWriter) {
	switch r.Format {
	case Format2.JSON:
		writeContentType(w, jsonContentType)
	case Format2.XML:
		writeContentType(w, xmlContentType)
	case Format2.CSV:
		writeContentType(w, csvContentType)
	}
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
