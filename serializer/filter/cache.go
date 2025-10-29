package filter

import (
	"reflect"
	"sort"
	"strings"
	"sync"
)

type cacheKey struct {
	typeName string
	groups   string
}

type fieldMapping struct {
	srcIndex int
	destIndex int
}

type typeCacheEntry struct {
	filteredType reflect.Type
	fieldMappings []fieldMapping
}

type typeCache struct {
	mu    sync.RWMutex
	cache map[cacheKey]*typeCacheEntry
}

var globalCache = &typeCache{
	cache: make(map[cacheKey]*typeCacheEntry),
}

func makeGroupKey(groups []string) string {
	if len(groups) == 0 {
		return ""
	}
	sorted := make([]string, len(groups))
	copy(sorted, groups)
	sort.Strings(sorted)
	return strings.Join(sorted, ",")
}

func makeCacheKey(t reflect.Type, groups []string) cacheKey {
	pkgPath := t.PkgPath()
	name := t.Name()
	var typeName string
	if pkgPath != "" {
		typeName = pkgPath + "." + name
	} else {
		typeName = name
	}

	return cacheKey{
		typeName: typeName,
		groups:   makeGroupKey(groups),
	}
}

func (tc *typeCache) Get(t reflect.Type, groups []string) (*typeCacheEntry, bool) {
	key := makeCacheKey(t, groups)
	tc.mu.RLock()
	entry, ok := tc.cache[key]
	tc.mu.RUnlock()
	return entry, ok
}

func (tc *typeCache) Set(t reflect.Type, groups []string, entry *typeCacheEntry) {
	key := makeCacheKey(t, groups)
	tc.mu.Lock()
	tc.cache[key] = entry
	tc.mu.Unlock()
}
