package api

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetMetaWithoutQueryParameters(t *testing.T) {
	assert := assert.New(t)

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.GET("/", MetaHandler)

	req, _ := http.NewRequest("GET", "/", nil)
	r.ServeHTTP(w, req)

	assert.Equal(
		http.StatusBadRequest,
		w.Code,
		"request without any query params should receive BadRequest",
	)
}
