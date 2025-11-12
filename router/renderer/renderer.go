package renderer

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
	jsonContentType        = []string{"application/json; charset=utf-8"}
	jsonldContentType      = []string{"application/ld+json; charset=utf-8"}
	xmlContentType         = []string{"application/xml; charset=utf-8"} // Standard RFC 7303
	csvContentType         = []string{"text/csv; charset=utf-8"}        // Standard RFC 4180
	messagepackContentType = []string{"application/msgpack"}
)

var serializerPools = map[format.Format]*sync.Pool{
	format.JSON: {
		New: func() any {
			return serializer.NewSerializer(format.JSON)
		},
	},
	format.JSONLD: {
		New: func() any {
			return serializer.NewSerializer(format.JSONLD)
		},
	},
	format.XML: {
		New: func() any {
			return serializer.NewSerializer(format.XML)
		},
	},
	format.CSV: {
		New: func() any {
			return serializer.NewSerializer(format.CSV)
		},
	},
	format.MESSAGEPACK: {
		New: func() any {
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

	str, err := s.Serialize(r.Data, serializer.SerializeOptions{
		Groups: r.Groups,
	})
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(str))
	return err
}

// WriteContentType sets the Content-Type header based on the format.
// It now uses GetContentType to centralize the logic.
func (r SerializerRenderer) WriteContentType(w http.ResponseWriter) {
	contentType := GetContentType(r.Format)
	if len(contentType) > 0 {
		// The writeContentType helper expects a []string.
		writeContentType(w, contentType)
	}
}

func writeContentType(w http.ResponseWriter, value []string) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// GetContentType returns the full content type string for a given format.
func GetContentType(f format.Format) []string {
	switch f {
	case format.JSON:
		return jsonContentType
	case format.JSONLD:
		return jsonldContentType
	case format.XML:
		return xmlContentType
	case format.CSV:
		return csvContentType
	case format.MESSAGEPACK:
		return messagepackContentType
	default:
		// Fallback to the format itself if unknown, though this should be rare.
		return []string{string(f)}
	}
}
