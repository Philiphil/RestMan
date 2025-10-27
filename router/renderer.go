package router

import (
	"net/http"
	"sync"

	"github.com/philiphil/restman/format"
	"github.com/philiphil/restman/serializer"
)

type SerializerRenderer struct {
	Data   any
	Format format.Format
	Groups []string
}

var (
	jsonContentType       = []string{"application/json; charset=utf-8"}
	jsonldContentType     = []string{"application/ld-json; charset=utf-8"}
	xmlContentType        = []string{"application/xml; charset=utf-8"}
	csvContentType        = []string{"text/csv"}
	messagepackContentType = []string{"application/msgpack"}
)

var serializerPools = map[format.Format]*sync.Pool{
	format.JSON: {
		New: func() interface{} {
			return serializer.NewSerializer(format.JSON)
		},
	},
	format.JSONLD: {
		New: func() interface{} {
			return serializer.NewSerializer(format.JSONLD)
		},
	},
	format.XML: {
		New: func() interface{} {
			return serializer.NewSerializer(format.XML)
		},
	},
	format.CSV: {
		New: func() interface{} {
			return serializer.NewSerializer(format.CSV)
		},
	},
	format.MESSAGEPACK: {
		New: func() interface{} {
			return serializer.NewSerializer(format.MESSAGEPACK)
		},
	},
}

func getSerializer(f format.Format) *serializer.Serializer {
	if pool, ok := serializerPools[f]; ok {
		return pool.Get().(*serializer.Serializer)
	}
	return serializer.NewSerializer(f)
}

func putSerializer(f format.Format, s *serializer.Serializer) {
	if pool, ok := serializerPools[f]; ok {
		pool.Put(s)
	}
}

// Render
func (r SerializerRenderer) Render(w http.ResponseWriter) (err error) {
	r.WriteContentType(w)
	s := getSerializer(r.Format)
	defer putSerializer(r.Format, s)

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
	case format.JSON:
		writeContentType(w, jsonContentType)
	case format.JSONLD:
		writeContentType(w, jsonldContentType)
	case format.XML:
		writeContentType(w, xmlContentType)
	case format.CSV:
		writeContentType(w, csvContentType)
	case format.MESSAGEPACK:
		writeContentType(w, messagepackContentType)
	}
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
