package middleware_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	app_mock "github.com/jonasjesusamerico/goexpert-rate-limiter/internal/application/mock"
	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/infra/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimiterMiddlewareTestSuite struct {
	suite.Suite
	mockRateLimiterApp *app_mock.RateLimiterServiceInterface
	middleware         *middleware.RateLimiterMiddleware
}

func (suite *RateLimiterMiddlewareTestSuite) SetupTest() {
	suite.mockRateLimiterApp = app_mock.NewRateLimiterServiceInterface(suite.T())
	suite.middleware = middleware.NewRateLimiterMiddleware(suite.mockRateLimiterApp)
}

func (suite *RateLimiterMiddlewareTestSuite) TestHandler_AllowRequest() {
	req := httptest.NewRequest("GET", "http://test.com/foo", nil)
	req.RemoteAddr = "192.168.0.1"
	req.Header.Set("API_KEY", "valid_token")

	recorder := httptest.NewRecorder()

	suite.mockRateLimiterApp.On("AllowRequest", mock.Anything, mock.Anything).Return(true)

	suite.middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})).ServeHTTP(recorder, req)

	assert.Equal(suite.T(), http.StatusOK, recorder.Code)
	assert.Equal(suite.T(), "OK", recorder.Body.String())
	suite.mockRateLimiterApp.AssertExpectations(suite.T())
}

func (suite *RateLimiterMiddlewareTestSuite) TestHandler_DenyRequest() {
	req := httptest.NewRequest("GET", "http://teste.com/foo", nil)
	req.RemoteAddr = "192.168.0.2"
	req.Header.Set("API_KEY", "invalid_token")

	recorder := httptest.NewRecorder()

	suite.mockRateLimiterApp.On("AllowRequest", mock.Anything, mock.Anything).Return(false)

	suite.middleware.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})).ServeHTTP(recorder, req)

	assert.Equal(suite.T(), http.StatusTooManyRequests, recorder.Code)
	suite.mockRateLimiterApp.AssertExpectations(suite.T())
}

func (suite *RateLimiterMiddlewareTestSuite) TestGetIP() {
	req1 := httptest.NewRequest("GET", "http://test.com", nil)
	req1.Header.Set("X-Forwarded-For", "192.168.0.1, 10.0.0.1, 172.16.0.1")
	ip1 := middleware.ExtractClientIP(req1)
	assert.Equal(suite.T(), "192.168.0.1", ip1, "Expected IP to be the first in X-Forwarded-For list")

	req2 := httptest.NewRequest("GET", "http://test.com", nil)
	req2.Header.Set("X-Forwarded-For", "192.168.0.1")
	ip2 := middleware.ExtractClientIP(req2)
	assert.Equal(suite.T(), "192.168.0.1", ip2, "Expected IP to be the single IP in X-Forwarded-For")

	req3 := httptest.NewRequest("GET", "http://test.com", nil)
	req3.RemoteAddr = "192.168.0.2:12345"
	ip3 := middleware.ExtractClientIP(req3)
	assert.Equal(suite.T(), "192.168.0.2", ip3, "Expected IP to be extracted from RemoteAddr")

	req4 := httptest.NewRequest("GET", "http://test.com", nil)
	req4.RemoteAddr = "invalid_address"
	ip4 := middleware.ExtractClientIP(req4)
	assert.Equal(suite.T(), "", ip4, "Expected empty IP due to invalid RemoteAddr format")
}

func TestRateLimiterMiddlewareTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterMiddlewareTestSuite))
}
