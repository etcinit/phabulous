package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jacobstr/confer"
	"github.com/stretchr/testify/assert"
)

func Test_EngineService_New(t *testing.T) {
	service := EngineService{}

	service.Config = confer.NewConfig()
	service.Config.Set("server.debug", false)

	engine := service.New()

	req, _ := http.NewRequest("GET", "/", nil)
	recorder := httptest.NewRecorder()

	engine.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
}

func Test_EngineService_New_WithDebug(t *testing.T) {
	service := EngineService{}

	service.Config = confer.NewConfig()
	service.Config.Set("server.debug", true)

	service.New()
}
