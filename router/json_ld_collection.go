package router

import (
	"strconv"
	"strings"
)

func JsonldCollection[T any](items []T, currentUrl string, currentPage int, params map[string]string, maxpage int) (m map[string]any) {
	m = map[string]any{}
	m["hydra:totalItems"] = len(items)
	//m["@context"] = "/api/contexts/" not handled
	m["hydra:member"] = items
	view := map[string]string{}
	view["@id"] = currentUrl
	view["@type"] = "hydra:PartialCollectionView"

	url := currentUrl[:strings.Index(currentUrl, "?")+1]
	m["@id"] = url[:len(url)-1]
	for k, v := range params {
		if k != "page" {
			url += k + "=" + v + "&"
		}
	}

	if currentPage != 1 { //dont bother to cast
		view["hydra:previous"] = url + "page=" + strconv.Itoa(currentPage-1)
	}
	if currentPage < maxpage {
		view["hydra:next"] = url + "page=" + strconv.Itoa(currentPage+1)
	}

	view["hydra:first"] = url + "page=1"
	view["hydra:last"] = url + "page=" + strconv.Itoa(max(maxpage, 1))

	m["hydra:view"] = view

	return m
}
