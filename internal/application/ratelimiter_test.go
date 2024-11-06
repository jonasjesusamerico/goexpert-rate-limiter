package application

import (
	"context"
	"testing"

	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/domain/service"

	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/domain/service/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimiterAppTestSuite struct {
	suite.Suite
	ctx                    context.Context
	mockRateLimiterService *mock.RateLimiterServiceInterface
	rateLimiterApp         *RateLimiterApp
}

func (suite *RateLimiterAppTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRateLimiterService = mock.NewRateLimiterServiceInterface(suite.T())
	suite.rateLimiterApp = NewRateLimiterApp(&service.RateLimiterService{RateLimiterRepo: suite.mockRateLimiterService})
}

func (suite *RateLimiterAppTestSuite) TestAllowRequest_AllowsRequest() {
	suite.mockRateLimiterService.On("AllowRequest", "127.0.0.1", "").Return(true)

	result := suite.rateLimiterApp.AllowRequest("127.0.0.1", "")
	suite.True(result)
	suite.mockRateLimiterService.AssertExpectations(suite.T())
}

func (suite *RateLimiterAppTestSuite) TestAllowRequest_DeniesRequest() {
	suite.mockRateLimiterService.On("AllowRequest", "127.0.0.1", "").Return(false)

	result := suite.rateLimiterApp.AllowRequest("127.0.0.1", "")
	suite.False(result)
	suite.mockRateLimiterService.AssertExpectations(suite.T())
}

func TestRateLimiterAppTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterAppTestSuite))
}
