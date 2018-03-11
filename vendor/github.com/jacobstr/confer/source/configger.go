package source

// Abstract various configuration sources into something that can get and set
// values, as well as return all of it's currently configured keys.
type Configger interface {
	// Get a value.
	Get(key string) (val interface{}, exists bool)
	// Set a value.
	Set(key string, val interface{})
	// Set data from a map[string]interface{}.
	FromStringMap(data map[string]interface{})
	// Merge a map[string]interface{} into existing data.
	ToStringMap() map[string]interface{}
}
