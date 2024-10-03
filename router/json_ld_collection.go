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

func Max(vars ...int) int {
	max := vars[0]

	for _, i := range vars {
		if i > max {
			max = i
		}
	}

	return max
}
