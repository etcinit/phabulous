package source

import (
	"reflect"
	"strings"

	"github.com/jacobstr/confer/maps"
	"github.com/spf13/cast"

	jww "github.com/spf13/jwalterweatherman"
)

// Manages key/value access for a specific configuration source. Delegated to by
// the over-arching config management functions that are aware of multiple config
// sources and their precedence.
type ConfigSource struct {
	// The raw configuration data.
	data map[string]interface{}

	// Hashmap of lower case keys to corresponding real key in data which we treat
	// as the canonical data store.
	index map[string]string
}

// Create a new case-insensitive, aliasable config map.
func NewConfigSource() *ConfigSource {
	return &ConfigSource{
		data:  make(map[string]interface{}),
		index: make(map[string]string),
	}
}

// Get the value at a key. Case-insensitive, but preserving.
func (self *ConfigSource) Get(key string) (val interface{}, exists bool) {
	index_key, index_exists := self.index[strings.ToLower(key)]

	// Exit if the index doesn't exist. We shouldn't have false negatives
	// unless our index falls out of sync.
	if index_exists == false {
		return nil, false
	}

	// Begin splitting the key apart.
	path := strings.Split(index_key, ".")
	current := self.data
	for _, part := range path[:len(path)-1] {
		if reflect.TypeOf(current).Kind() != reflect.Map {
			jww.TRACE.Println("Attempting deep access of a non-map.")
			return nil, false
		} else {
			var next interface{}
			next, exists := current[part]
			if exists == false {
				return nil, false
			} else {
				current = cast.ToStringMap(next)
			}
		}
	}

	val, exists = current[path[len(path)-1]]

	// Use a helper function if one is provided.
	switch v := val.(type) {
	case func() interface{}:
		return v(), exists
	default:
		return v, exists
	}
}

// Set a key in a case insensitive manner.
func (self *ConfigSource) Set(key string, val interface{}) {
	index_key, index_exists := self.index[strings.ToLower(key)]
	if index_exists == false {
		index_key = key
	}

	path := strings.Split(index_key, ".")
	original_path := strings.Split(key, ".")

	current := self.data
	for depth, part := range path[:len(path)-1] {
		if reflect.TypeOf(current).Kind() != reflect.Map {
			panic("Attempting deep access of a non-map.")
		} else {
			var next interface{}
			next, exists := current[part]

			// Generate the index of our ancestors as we progress.
			ancestor_key := strings.Join(original_path[0:depth+1], ".")

			// Stub out ancestors if we're setting a deep child.
			if exists == false {
				current[part] = make(map[string]interface{})
				current = current[part].(map[string]interface{})
				self.index[strings.ToLower(ancestor_key)] = ancestor_key
			} else {
				current = next.(map[string]interface{})
			}
		}
	}

	current[path[len(path)-1]] = val
	self.updateIndex(key, current)
}

// Replaces our configuration data with the provided stringmap, without merging.
func (self *ConfigSource) FromStringMap(data map[string]interface{}) {
	self.data = data
	self.UpdateIndices()
}

// Returns data as a string map.
func (self *ConfigSource) ToStringMap() map[string]interface{} {
	return self.data
}

// Updates our lookup table of insensitive materialized paths to their
// corresponding 'real' keys. E.g.
//
//		Database.Connections.Hosts <- database.connections.hosts
//
// By maintaining a separate index and maintaining case in the original
// stringmaps (e.g. by lowercasing keys directly) we accomodate the passing
// of config data to structures that ~may~ be case sensitive. I.E we avoid
// destructive operations on configurationd data.
func (self *ConfigSource) updateIndex(key string, data interface{}) {
	if data == nil {
		return
	}

	// Don't change the case of the original key if it already exists.
	_, index_exists := self.index[strings.ToLower(key)]
	if index_exists == false {
		self.index[strings.ToLower(key)] = key
	}

	if reflect.TypeOf(data).Kind() != reflect.Map {
		return
	}

	for child_key, val := range cast.ToStringMap(data) {
		var joined_key string
		if len(key) > 0 {
			joined_key = key + "." + child_key
		} else {
			joined_key = child_key
		}
		self.updateIndex(joined_key, val)
	}
}

// Index every key/value pair inside of this config sources's data.
func (self *ConfigSource) UpdateIndices() {
	for key, val := range self.data {
		jww.TRACE.Println("update index", key)
		self.updateIndex(key, val)
	}
}

// Returns all the keys for this specific configuration source.
func (self *ConfigSource) AllKeys() []string {
	return maps.CollectKeys(self.data, "", -1)
}
