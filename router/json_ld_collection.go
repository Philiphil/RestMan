package router

import (
	"strconv"
	"strings"
)

// Max returns the maximum value from the provided integers.
func Max(vars ...int) int {
	max := vars[0]

	for _, i := range vars {
		if i > max {
			max = i
		}
	}

	return max
}

// HydraCollection represents a fully Hydra-compliant JSON-LD collection structure.
type HydraCollection struct {
	Context    any          `json:"@context,omitempty"`
	ID         string       `json:"@id"`
	Type       string       `json:"@type"`
	TotalItems int          `json:"hydra:totalItems"`
	Member     any          `json:"hydra:member"`
	View       *HydraView   `json:"hydra:view,omitempty"`
	Search     *HydraSearch `json:"hydra:search,omitempty"`
}

// HydraView represents the pagination view within a Hydra collection.
type HydraView struct {
	ID       string `json:"@id"`
	Type     string `json:"@type"`
	First    string `json:"hydra:first,omitempty"`
	Last     string `json:"hydra:last,omitempty"`
	Previous string `json:"hydra:previous,omitempty"`
	Next     string `json:"hydra:next,omitempty"`
}

// HydraSearch represents the search template for filtering collections.
type HydraSearch struct {
	Type     string             `json:"@type"`
	Template string             `json:"hydra:template"`
	Mapping  []HydraIriTemplate `json:"hydra:mapping,omitempty"`
}

// HydraIriTemplate represents a single mapping in the search template.
type HydraIriTemplate struct {
	Type     string `json:"@type"`
	Variable string `json:"hydra:variable"`
	Property string `json:"hydra:property,omitempty"`
	Required bool   `json:"hydra:required"`
}

// HydraOperation represents an available operation on a resource.
type HydraOperation struct {
	Type        string `json:"@type"`
	Method      string `json:"hydra:method"`
	Label       string `json:"hydra:title,omitempty"`
	Description string `json:"hydra:description,omitempty"`
	Expects     string `json:"hydra:expects,omitempty"`
	Returns     string `json:"hydra:returns,omitempty"`
	StatusCodes []int  `json:"hydra:possibleStatus,omitempty"`
}

// HydraResource represents a JSON-LD resource with metadata.
type HydraResource map[string]any

// JsonldCollection creates a Hydra-compliant JSON-LD collection with pagination metadata.
// totalItems should be the total count across all pages (for proper Hydra spec compliance).
func JsonldCollection[T any](items []T, totalItems int, currentUrl string, currentPage int, params map[string]string, maxpage int) (m map[string]any) {
	m = map[string]any{}
	m["hydra:totalItems"] = totalItems // Total count across all pages, not just len(items)
	//m["@context"] = "/api/contexts/" not handled
	m["hydra:member"] = items
	view := map[string]string{}
	view["@id"] = currentUrl
	view["@type"] = "hydra:PartialCollectionView"

	url := currentUrl[:strings.Index(currentUrl, "?")+1]
	m["@id"] = url[:len(url)-1]
	stringifyParam := ""
	for k, v := range params {
		if k != "page" {
			stringifyParam += k + "=" + v + "&"
		}
	}
	url += stringifyParam

	if currentPage != 1 { //dont bother to cast
		view["hydra:previous"] = url + "page=" + strconv.Itoa(currentPage-1)
	}
	if currentPage < maxpage {
		view["hydra:next"] = url + "page=" + strconv.Itoa(currentPage+1)
	}

	view["hydra:first"] = url + "page=1"
	view["hydra:last"] = url + "page=" + strconv.Itoa(Max(maxpage, 1))

	m["hydra:view"] = view

	return m
}
