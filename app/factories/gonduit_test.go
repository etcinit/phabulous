package factories

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/jacobstr/confer"
	"github.com/stretchr/testify/assert"
)

func getConfig() *confer.Config {
	c := confer.NewConfig()

	c.Set("mist.ignore-ca", false)

	return c
}

func getServer(failGetCaps bool, failConnect bool) *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	r.POST("/api/conduit.getcapabilities", func(c *gin.Context) {
		if failGetCaps {
			c.String(400, "failed")
		}

		c.JSON(200, gin.H{
			"result": map[string][]string{
				"authentication": []string{"token", "session"},
				"signatures":     []string{"consign"},
				"input":          []string{"json", "urlencoded"},
				"output":         []string{"json"},
			},
		})
	})

	r.POST("/api/conduit.connect", func(c *gin.Context) {
		if failConnect {
			c.String(400, "failed")
		}

		c.JSON(200, gin.H{
			"result": gin.H{
				"connectionID": 23,
				"sessionKey":   "some key",
			},
		})
	})

	return r
}

func TestMake(t *testing.T) {
	ts := httptest.NewServer(getServer(false, false))
	defer ts.Close()

	c := getConfig()
	c.Set("conduit.api", ts.URL)
	c.Set("conduit.token", "api-sometoken")

	factory := &GonduitFactory{Config: c}
	client, err := factory.Make()

	assert.NotNil(t, client)
	assert.Nil(t, err)
}

func TestMake_withError(t *testing.T) {
	ts := httptest.NewServer(getServer(true, false))
	defer ts.Close()

	c := getConfig()
	c.Set("conduit.api", ts.URL)
	c.Set("conduit.token", "api-sometoken")

	factory := &GonduitFactory{Config: c}
	client, err := factory.Make()

	assert.Nil(t, client)
	assert.NotNil(t, err)
}

func TestMake_withCertificate(t *testing.T) {
	ts := httptest.NewServer(getServer(false, false))
	defer ts.Close()

	c := getConfig()
	c.Set("conduit.api", ts.URL)
	c.Set("conduit.user", "some-user")
	c.Set("conduit.cert", "some-cert")

	factory := &GonduitFactory{Config: c}
	client, err := factory.Make()

	assert.NotNil(t, client)
	assert.Nil(t, err)
}

func TestMake_withCertificateError(t *testing.T) {
	ts := httptest.NewServer(getServer(false, true))
	defer ts.Close()

	c := getConfig()
	c.Set("conduit.api", ts.URL)
	c.Set("conduit.user", "some-user")
	c.Set("conduit.cert", "some-cert")

	factory := &GonduitFactory{Config: c}
	client, err := factory.Make()

	assert.Nil(t, client)
	assert.NotNil(t, err)
}

func TestMake_withMissingCredentials(t *testing.T) {
	ts := httptest.NewServer(getServer(false, false))
	defer ts.Close()

	c := getConfig()
	c.Set("conduit.api", ts.URL)

	factory := &GonduitFactory{Config: c}
	client, err := factory.Make()

	assert.Nil(t, client)
	assert.Equal(t, ErrNoCredentials, err)
}
