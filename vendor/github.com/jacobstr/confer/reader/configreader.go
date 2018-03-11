package reader

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/jacobstr/confer/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gopkg.in/yaml.v1"
)

type ConfigFormat string

const (
	FormatYAML ConfigFormat = "yaml"
	FormatJSON ConfigFormat = "json"
	FormatTOML ConfigFormat = "toml"
)

type ConfigReader struct {
	Format string
	reader io.Reader
}

// Retuns the configuration data into a generic object for for us.
func (cr *ConfigReader) Export() (interface{}, error) {
	var config interface{}
	buf := new(bytes.Buffer)
	buf.ReadFrom(cr.reader)

	switch cr.Format {
	case "yaml":
		if err := yaml.Unmarshal(buf.Bytes(), &config); err != nil {
			jww.ERROR.Fatalf("Error parsing config: %s", err)
		}

	case "json":
		if err := json.Unmarshal(buf.Bytes(), &config); err != nil {
			jww.ERROR.Fatalf("Error parsing config: %s", err)
		}

	case "toml":
		if _, err := toml.Decode(buf.String(), &config); err != nil {
			jww.ERROR.Fatalf("Error parsing config: %s", err)
		}
	default:
		return nil, err.UnsupportedConfigError(cr.Format)
	}

	return config, nil
}

func ReadFile(path string) (interface{}, error) {
	file, err := ioutil.ReadFile(path)
	if err != nil {
		jww.DEBUG.Println("Error reading config file:", err)
		return nil, err
	}

	reader := bytes.NewReader(file)

	cr := &ConfigReader{Format: getConfigType(path), reader: reader}
	return cr.Export()
}

func ReadBytes(data []byte, format string) (interface{}, error) {
	cr := ConfigReader{
		Format: format,
		reader: bytes.NewReader(data),
	}

	return cr.Export()
}

func getConfigType(path string) string {
	ext := filepath.Ext(path)
	switch ext[1:] {
	case "yml":
		return "yaml"
	default:
		return ext[1:]
	}
}
