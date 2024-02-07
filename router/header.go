package router

import (
	"github.com/philiphil/apiman/errors"
	"github.com/philiphil/apiman/format"
	"sort"
	"strconv"
	"strings"
)

type MediaType struct {
	Type   string
	Weight float64
}

func ParseAcceptHeader(acceptHeader string) (format.Format, error) {
	if acceptHeader == "" {
		return format.JSON, nil
	}
	mediaTypes := strings.Split(acceptHeader, ",")
	mediaTypesWithQ := make([]MediaType, len(mediaTypes))

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

	sortedMediaTypes := make([]string, len(mediaTypesWithQ))
	for i, mediaType := range mediaTypesWithQ {
		sortedMediaTypes[i] = mediaType.Type
		if f := ParseTypeFromString(mediaType.Type); f != format.Undefined && f != format.Unknown {
			return f, nil
		} else if mediaType.Type == "*/*" {
			return format.JSON, nil
		}
	} //default
	return format.Undefined, errors.ErrNotAcceptable
}

func ParseTypeFromString(str string) format.Format {
	if str == "" {
		return format.Undefined
	}
	if strings.Contains(strings.ToLower(str), "json") {
		return format.JSON
	}
	if strings.Contains(strings.ToLower(str), "xml") {
		return format.XML
	}
	if strings.Contains(strings.ToLower(str), "csv") {
		return format.CSV
	}
	return format.Unknown
}
