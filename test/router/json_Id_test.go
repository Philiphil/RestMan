package router_test

import (
	"testing"

	"reflect"

	. "github.com/philiphil/restman/router"
)

func TestMax(t *testing.T) {
	if Max(1, 2, 3, 4, 5) != 5 {
		t.Error("!")
	}
}

func TestJsonldCollection(t *testing.T) {
	items := []int{1, 2, 3, 4, 5}
	currentUrl := "http://localhost:8080/api/test?pagination=true&page=2"
	currentPage := 2
	params := map[string]string{"page": "2", "pagination": "true"}
	maxpage := 3
	m := JsonldCollection(items, currentUrl, currentPage, params, maxpage)
	if m["hydra:totalItems"] != 5 {
		t.Error("!")
	}
	if !reflect.DeepEqual(m["hydra:member"], items) {
		t.Error("!")
	}
	view := m["hydra:view"].(map[string]string)
	if view["@id"] != "http://localhost:8080/api/test?pagination=true&page=2" {
		t.Error(view["@id"])
	}
	if view["@type"] != "hydra:PartialCollectionView" {
		t.Error("!")
	}
	if view["hydra:previous"] != "http://localhost:8080/api/test?pagination=true&page=1" {
		t.Error("!")
	}
	if view["hydra:next"] != "http://localhost:8080/api/test?pagination=true&page=3" {
		t.Error("!")
	}
	if view["hydra:first"] != "http://localhost:8080/api/test?pagination=true&page=1" {
		t.Error("!")
	}
	if view["hydra:last"] != "http://localhost:8080/api/test?pagination=true&page=3" {
		t.Error("!")
	}
}
