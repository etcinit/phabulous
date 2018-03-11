package source

import (
	"fmt"
	"os"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

// A configuration data source that that reads environment variables.
type EnvSource struct {
	index map[string]string
}

// Converts our materialized path format to a corresponding ENV_VAR friendly
// format. Periods are replaced with single underscores. Note that reversing
// this would generally be ambiguous as underscores are common in variable keys.
func envamize(key string) string {
	return strings.Replace(strings.ToUpper(key), ".", "_", -1)
}

func NewEnvSource() *EnvSource {
	return &EnvSource{
		index: make(map[string]string),
	}
}

// Essentially an environment variable specific alias.
func (self *EnvSource) Bind(input ...string) (err error) {
	var key, envkey string

	if len(input) == 0 {
		return fmt.Errorf("BindEnv missing key to bind to")
	}

	if len(input) == 1 {
		key = input[0]
	} else {
		key = input[1]
	}

	envkey = envamize(key)

	jww.TRACE.Println(key, "Bound to", envkey)
	self.index[strings.ToLower(key)] = envkey

	return nil
}

func (self *EnvSource) AllKeys() []string {
	a := []string{}
	for x, _ := range self.index {
		a = append(a, strings.ToLower(x))
	}
	return a
}

// Gets an environment variable.
func (self *EnvSource) Get(key string) (val interface{}, exists bool) {
	envkey, exists := self.index[key]

	if exists {
		jww.TRACE.Println(key, "registered as env var", envkey)
	}

	if val = os.Getenv(envkey); val != "" {
		jww.TRACE.Println(envkey, "found in environment with val:", val)
		return val, true
	} else {
		jww.TRACE.Println(envkey, "env value unset:")
		return nil, false
	}
}
