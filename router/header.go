package router

import (
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/philiphil/restman/errors"
	"github.com/philiphil/restman/format"
)

// MediaType represents a media type with its quality weight from the Accept header.
type MediaType struct {
	Type   string
	Weight float64
}

var (
	acceptHeaderCache sync.Map
	mediaTypePool     = sync.Pool{
		New: func() interface{} {
			return make([]MediaType, 0, 8)
		},
	}
)

// ParseAcceptHeader parses the Accept HTTP header and returns the most preferred supported format.
func ParseAcceptHeader(acceptHeader string) (format.Format, error) {
	if acceptHeader == "" {
		return format.JSON, nil
	}

	if cached, ok := acceptHeaderCache.Load(acceptHeader); ok {
		return cached.(format.Format), nil
	}

	mediaTypes := strings.Split(acceptHeader, ",")
	mediaTypesWithQ := mediaTypePool.Get().([]MediaType)
	defer func() {
		mediaTypesWithQ = mediaTypesWithQ[:0]
		mediaTypePool.Put(mediaTypesWithQ)
	}()

	if cap(mediaTypesWithQ) < len(mediaTypes) {
		mediaTypesWithQ = make([]MediaType, len(mediaTypes))
	} else {
		mediaTypesWithQ = mediaTypesWithQ[:len(mediaTypes)]
	}

	for i, mediaType := range mediaTypes {
		parts := strings.Split(strings.TrimSpace(mediaType), ";")
		mediaTypesWithQ[i].Type = parts[0]

		mediaTypesWithQ[i].Weight = 1.0
		for _, part := range parts[1:] {
			if pos := strings.Index(part, "q="); pos > -1 {
				qValue, err := strconv.ParseFloat(strings.TrimSpace(part[pos+2:]), 64)
				if err == nil {
					mediaTypesWithQ[i].Weight = qValue
				}
			}
		}
	}

	sort.Slice(mediaTypesWithQ, func(i, j int) bool {
		if mediaTypesWithQ[i].Weight == mediaTypesWithQ[j].Weight {
			return mediaTypesWithQ[i].Type < mediaTypesWithQ[j].Type
		}
		return mediaTypesWithQ[i].Weight > mediaTypesWithQ[j].Weight
	})

	for _, mediaType := range mediaTypesWithQ {
		if f := ParseTypeFromString(mediaType.Type); f != format.Undefined && f != format.Unknown {
			acceptHeaderCache.Store(acceptHeader, f)
			return f, nil
		} else if mediaType.Type == "*/*" {
			acceptHeaderCache.Store(acceptHeader, format.JSON)
			return format.JSON, nil
		}
	}

	return format.Undefined, errors.ErrNotAcceptable
}

// ParseTypeFromString converts a media type string to a Format constant.
func ParseTypeFromString(str string) format.Format {
	if str == "" {
		return format.Undefined
	}
	lower := strings.ToLower(str)
	if strings.Contains(lower, "ld+json") {
		return format.JSONLD
	}
	if strings.Contains(lower, "json") {
		return format.JSON
	}
	if strings.Contains(lower, "xml") {
		return format.XML
	}
	if strings.Contains(lower, "csv") {
		return format.CSV
	}
	return format.Unknown
}
