// Copyright © 2014 Steve Francia <spf@spf13.com>.
// Copyright © 2014 Jacob Straszysnki <jacobstr@gmail.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Confer is a application configuration system.
// It believes that applications can be configured a variety of ways
// via flags, ENVIRONMENT variables, configuration files retrieved
// from the file system.
//
// There are 3 precedence tiers:
//
// 1. Command line flags.
// 2. Environment variables.
// 3. Attributes - (e.g. Set, SetDefault, ReadPaths)

package confer

import (
	"fmt"
	"path"
	"path/filepath"
	"strings"
	"time"
	"reflect"

	"github.com/kr/pretty"
	"github.com/spf13/cast"
	jww "github.com/spf13/jwalterweatherman"
	"github.com/spf13/pflag"

	"github.com/jacobstr/confer/reader"
	. "github.com/jacobstr/confer/source"

	errors "github.com/jacobstr/confer/errors"
	"github.com/jacobstr/confer/maps"
)

// Manages key/value access and aliasing across multiple configuration sources.
type Config struct {
	pflags     *PFlagSource
	env        *EnvSource
	attributes *ConfigSource

	// The root path for configuration files.
	rootPath string
}

func NewConfig() *Config {
	manager := &Config{}
	manager.pflags = NewPFlagSource()
	manager.attributes = NewConfigSource()
	manager.env = NewEnvSource()
	manager.rootPath = ""

	return manager
}

// Finds a value at a provided key, returning nil if the key does not exist.
// The order of precedence for configuration data is:
// 1. Program arguments.
// 2. Environment variables.
// 3. Config file data, overrides, and defaults.
func (self *Config) Find(key string) interface{} {
	var val interface{}
	var exists bool

	// PFlag Override first
	val, exists = self.pflags.Get(key)
	if exists {
		jww.TRACE.Println(key, "found in override (via pflag):", val)
		return val
	}

	// Periods are not supported. Allow the usage of underscores to specify nested
	// configuration options.
	val, exists = self.env.Get(key)
	if exists {
		jww.TRACE.Println(key, "Found in environment with value:", val)
		return val
	}

	// Attributes entail pretty much everything else.
	val, exists = self.attributes.Get(key)
	if exists {
		jww.TRACE.Println(key, "Found in config:", val)
		return val
	}

	return nil
}

func (manager *Config) GetString(key string) string {
	return cast.ToString(manager.Get(key))
}

func (manager *Config) GetBool(key string) bool {
	return cast.ToBool(manager.Get(key))
}

func (manager *Config) GetInt(key string) int {
	return cast.ToInt(manager.Get(key))
}

func (manager *Config) GetFloat64(key string) float64 {
	return cast.ToFloat64(manager.Get(key))
}

func (manager *Config) GetTime(key string) time.Time {
	return cast.ToTime(manager.Get(key))
}

func (manager *Config) GetStringSlice(key string) []string {
	return cast.ToStringSlice(manager.Get(key))
}

func (manager *Config) GetStringMap(key string) map[string]interface{} {
	return cast.ToStringMap(manager.Get(key))
}

func (manager *Config) GetStringMapString(key string) map[string]string {
	return cast.ToStringMapString(manager.Get(key))
}

// Binds a configuration key to a command line flag:
//	 pflag.Int("port", 8080, "The best alternative port")
//	 confer.BindPFlag("port", pflag.Lookup("port"))
func (manager *Config) BindPFlag(key string, flag *pflag.Flag) (err error) {
	if flag == nil {
		return fmt.Errorf("flag for %q is nil", key)
	}

	manager.pflags.Set(key, flag)

	switch flag.Value.Type() {
	case "int", "int8", "int16", "int32", "int64":
		manager.SetDefault(key, cast.ToInt(flag.Value.String()))
	case "bool":
		manager.SetDefault(key, cast.ToBool(flag.Value.String()))
	default:
		manager.SetDefault(key, flag.Value.String())
	}
	return nil
}

// Binds a confer key to a ENV variable. ENV variables are case sensitive If only
func (manager *Config) BindEnv(input ...string) (err error) {
	return manager.env.Bind(input...)
}

// Get returns an interface..
// Must be typecast or used by something that will typecast
func (manager *Config) Get(key string) interface{} {
	jww.TRACE.Println("Looking for", key)

	v := manager.Find(key)

	if v == nil {
		return nil
	}

	jww.TRACE.Println("Found value", v)
	switch v.(type) {
	case bool:
		return cast.ToBool(v)
	case string:
		return cast.ToString(v)
	case int64, int32, int16, int8, int:
		return cast.ToInt(v)
	case float64, float32:
		return cast.ToFloat64(v)
	case time.Time:
		return cast.ToTime(v)
	case []string:
		return v
	}
	return v
}

// Returns true if the config key exists and is non-nil.
func (manager *Config) IsSet(key string) bool {
	t := manager.Get(key)
	return t != nil
}

// Have confer check ENV variables for all
// keys set in config, default & flags
func (manager *Config) AutomaticEnv() {
	for _, x := range manager.AllKeys() {
		manager.BindEnv(x)
	}
}

// Returns true if the key provided exists in our configuration.
func (manager *Config) InConfig(key string) bool {
	_, exists := manager.attributes.Get(key)
	return exists
}

// Set the default value for this key.
// Default only used when no value is provided by the user via flag, config or ENV.
func (manager *Config) SetDefault(key string, value interface{}) {
	if !manager.IsSet(key) {
		manager.attributes.Set(key, value)
	}
}

// Explicitly sets a value. Will not override command line arguments or
// environment variables, as those sources have higher precedence.
func (manager *Config) Set(key string, value interface{}) {
	manager.attributes.Set(key, value)
}

// Sets an optional root path. This frees you from having to specify a
// redundant prefix when calling ReadPaths() later.
func (manager *Config) SetRootPath(path string) {
	manager.rootPath = path
}

// Loads and sequentially + recursively merges the provided config arguments. Returns
// an error if any of the files fail to load, though this may be expecte
// in the case of search paths.
func (manager *Config) ReadPaths(paths ...string) error {
	var err error
	var loaded interface{}

	merged_config := manager.attributes.ToStringMap()
	errs := []error{}

	for _, base_path := range paths {
		var final_path string

		if filepath.IsAbs(base_path) == false {
			final_path = path.Join(manager.rootPath, base_path)
		} else {
			final_path = path.Join(manager.rootPath, base_path)
		}

		loaded, err = reader.ReadFile(final_path)

		if err != nil {
			errs = append(errs, err)
			continue
		}

		// In-place recursive coercion to stringmap.
		coerced := cast.ToStringMap(loaded)
		maps.ToStringMapRecursive(coerced)

		if merged_config == nil {
			merged_config = coerced
		} else {
			merged_config = maps.Merge(
				merged_config,
				coerced,
			)
		}

		manager.attributes.FromStringMap(merged_config)
	}

	if len(errs) > 0 {
		return &errors.LoadError{Errors: errs}
	} else {
		return nil
	}
}

// Merges data into the our attributes configuration tier from a struct.
func (manager *Config) MergeAttributes(val interface{}) error {
	merged_config := maps.Merge(
		manager.attributes.ToStringMap(),
		cast.ToStringMap(val),
	)

	manager.attributes.FromStringMap(merged_config)
	return nil
}

// Returns all currently set keys, pruning ancestors and only
// showing the leaves.
func (manager *Config) AllKeys() []string {
	keys := manager.attributes.AllKeys()
	keys = append(keys, manager.env.AllKeys()...)
	keys = append(keys, manager.attributes.AllKeys()...)

	leaves := map[string]struct{}{}
	for _, key := range keys {

		// Filter out leaves. This is really ineffecient.
		val := manager.Get(key)
		if val == nil {
			leaves[key] = struct{}{}
		} else if reflect.TypeOf(val).Kind() != reflect.Map {
			leaves[key] = struct{}{}
		}
	}

	unique_keys := []string{}
	for x, _ := range leaves {
		// LowerCase the key for backwards-compatibility.
		unique_keys = append(unique_keys, strings.ToLower(x))
	}

	return unique_keys
}

func (manager *Config) AllSettings() map[string]interface{} {
	m := map[string]interface{}{}
	for _, x := range manager.AllKeys() {
		m[x] = manager.Get(x)
	}

	return m
}

func (manager *Config) Debug() {
	fmt.Println("Flags:")
	pretty.Println(manager.pflags)
	fmt.Println("Env:")
	pretty.Println(manager.env)
	fmt.Println("Config file attributes:")
	pretty.Println(manager.attributes)
}
