package maps

import (
	"reflect"

	"github.com/spf13/cast"
)

// A callback for Traverse. It will accept a key, a value, the current depth
// and return true if we should continue deeper.
type Traverser func(key string, val interface{}, depth int) bool

// General purpose method for traversing a string map.
func Traverse(data map[string]interface{}, cb Traverser) {
	traverse(data, "", 0, cb)
}

// Generic functional, recursive stringmap traversal.
// Provides the callback with the current value, materialized path, and depth.
func traverse(data map[string]interface{}, path string, depth int, cb Traverser) {
	for key, val := range data {
		var joined_key string
		if len(path) > 0 {
			joined_key = path + "." + key
		} else {
			joined_key = key
		}

		if cb(joined_key, val, depth) {
			if val != nil && reflect.TypeOf(val).Kind() == reflect.Map {
				traverse(cast.ToStringMap(val), joined_key, depth+1, cb)
			}
		}
	}
}

// Recursively collects all keys into a flattened slice of materialized paths.
func CollectKeys(data map[string]interface{}, path string, max_depth int) []string {
	m := []string{}
	Traverse(data, func(key string, val interface{}, depth int) bool {
		m = append(m, key)
		return max_depth == -1 || depth <= max_depth
	})
	return m
}

// Adapted from github.com/peterbourgon/mergemap
var (
	MaxDepth = 32
)

// Merge recursively merges the src and dst maps. Key conflicts are resolved by
// preferring src, or recursively descending, if both src and dst are maps.
func Merge(dst, src map[string]interface{}) map[string]interface{} {
	return merge(dst, src, 0)
}

func merge(dst, src map[string]interface{}, depth int) map[string]interface{} {
	if depth > MaxDepth {
		panic("too deep!")
	}
	for key, srcVal := range src {
		if dstVal, ok := dst[key]; ok {
			srcMap, srcMapOk := mapify(srcVal)
			dstMap, dstMapOk := mapify(dstVal)
			if srcMapOk && dstMapOk {
				srcVal = merge(dstMap, srcMap, depth+1)
			}
		}
		dst[key] = srcVal
	}
	return dst
}

func mapify(i interface{}) (map[string]interface{}, bool) {
	v, err := cast.ToStringMapE(i)
	if err != nil {
		return v, false
	} else {
		return v, true
	}
}

// Recursively coerces all maps to a stringmap. Because that's how we want it.
func ToStringMapRecursive(src map[string]interface{}) {
	for key, val := range src {
		coerced, err := cast.ToStringMapE(val)
		if err == nil {
			src[key] = coerced
			ToStringMapRecursive(coerced)
		}
	}
}
