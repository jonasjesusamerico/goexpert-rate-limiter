package service

import (
	"context"
	"testing"

	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/domain/service/mock"
	"github.com/stretchr/testify/suite"
)

type RateLimiterServiceTestSuite struct {
	suite.Suite
	ctx                 context.Context
	mockRateLimiterRepo *mock.RateLimiterServiceInterface
	rateLimiterService  *RateLimiterService
}

func (suite *RateLimiterServiceTestSuite) SetupTest() {
	suite.ctx = context.Background()
	suite.mockRateLimiterRepo = mock.NewRateLimiterServiceInterface(suite.T())
	suite.rateLimiterService = NewRateLimiterService(suite.mockRateLimiterRepo)
}

func (suite *RateLimiterServiceTestSuite) TestAllowRequest_AllowsRequest() {
	suite.mockRateLimiterRepo.On("AllowRequest", "127.0.0.1", "").Return(true)

	result := suite.rateLimiterService.AllowRequest("127.0.0.1", "")
	suite.True(result)
	suite.mockRateLimiterRepo.AssertExpectations(suite.T())
}

func (suite *RateLimiterServiceTestSuite) TestAllowRequest_DeniesRequest() {
	suite.mockRateLimiterRepo.On("AllowRequest", "127.0.0.1", "").Return(false)

	result := suite.rateLimiterService.AllowRequest("127.0.0.1", "")
	suite.False(result)
	suite.mockRateLimiterRepo.AssertExpectations(suite.T())
}

func TestRateLimiterServiceTestSuite(t *testing.T) {
	suite.Run(t, new(RateLimiterServiceTestSuite))
}
