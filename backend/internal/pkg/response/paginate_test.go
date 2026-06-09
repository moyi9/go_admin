package response

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"go_web/internal/pkg/apperror"
)

func TestPageDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?a=1", nil)

	offset, limit := Page(c)
	require.Equal(t, 0, offset)
	require.Equal(t, DefaultPageSize, limit)
}

func TestPageCustomValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=3&page_size=10", nil)

	offset, limit := Page(c)
	require.Equal(t, 20, offset) // (3-1)*10
	require.Equal(t, 10, limit)
}

func TestPageCapsMaxSize(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=1&page_size=9999", nil)

	_, limit := Page(c)
	require.Equal(t, MaxPageSize, limit)
}

func TestPageIgnoresInvalidValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=abc&page_size=-1", nil)

	offset, limit := Page(c)
	require.Equal(t, 0, offset)        // page defaults to 1
	require.Equal(t, DefaultPageSize, limit) // page_size defaults
}

func TestPageIgnoresZeroValues(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())
	c.Request = httptest.NewRequest(http.MethodGet, "/?page=0&page_size=0", nil)

	offset, limit := Page(c)
	require.Equal(t, 0, offset)
	require.Equal(t, DefaultPageSize, limit)
}

func TestPaginatedOKResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("request_id", "req-123")

	PaginatedOK(c, []string{"a", "b"}, int64(42), 10, 10)

	require.Equal(t, http.StatusOK, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, `"code":0`)
	require.Contains(t, body, `"total":42`)
	require.Contains(t, body, `"page":2`)       // offset=10, limit=10 → page 2
	require.Contains(t, body, `"page_size":10`)
	require.Contains(t, body, `"request_id":"req-123"`)
}

func TestPaginatedResponsePageOne(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)

	PaginatedOK(c, nil, int64(0), 0, 20)

	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), `"page":1`)
}

func TestErrorResponseIncludesRequestID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	c.Request = httptest.NewRequest(http.MethodGet, "/", nil)
	c.Set("request_id", "err-req-id")

	Error(c, http.StatusBadRequest, apperror.CodeInvalidArgument, "bad input")

	require.Equal(t, http.StatusBadRequest, rec.Code)
	body := rec.Body.String()
	require.Contains(t, body, `"code":40000`)
	require.Contains(t, body, `"request_id":"err-req-id"`)
}
