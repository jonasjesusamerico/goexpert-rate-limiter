package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) SetupTest() {
	os.Setenv("MAX_REQUESTS_PER_SECOND", "200")
	os.Setenv("BLOCK_DURATION_SECONDS", "120")
	os.Setenv("IP_MAX_REQUESTS_PER_SECOND", "10")
	os.Setenv("IP_BLOCK_DURATION_SECONDS", "600")
	os.Setenv("TOKEN_MAX_REQUESTS_PER_SECOND", "20")
	os.Setenv("TOKEN_BLOCK_DURATION_SECONDS", "600")
	os.Setenv("REDIS_HOST", "testhost")
	os.Setenv("REDIS_PORT", "6380")
	os.Setenv("REDIS_PASSWORD", "testpassword")
}

func (suite *ConfigTestSuite) TearDownTest() {
	os.Unsetenv("MAX_REQUESTS_PER_SECOND")
	os.Unsetenv("BLOCK_DURATION_SECONDS")
	os.Unsetenv("IP_MAX_REQUESTS_PER_SECOND")
	os.Unsetenv("IP_BLOCK_DURATION_SECONDS")
	os.Unsetenv("TOKEN_MAX_REQUESTS_PER_SECOND")
	os.Unsetenv("TOKEN_BLOCK_DURATION_SECONDS")
	os.Unsetenv("REDIS_HOST")
	os.Unsetenv("REDIS_PORT")
	os.Unsetenv("REDIS_PASSWORD")
}

func (suite *ConfigTestSuite) TestLoadConfig() {
	cfg, err := LoadConfig(".env")
	suite.NoError(err)
	suite.Equal(10, cfg.IPMaxRequestsPerSecond)
	suite.Equal(120, cfg.BlockDurationSeconds)
	suite.Equal(600, cfg.IPBlockDurationSeconds)
	suite.Equal(600, cfg.TokenBlockDurationSeconds)
	suite.Equal(200, cfg.MaxRequestsPerSecond)
	suite.Equal(20, cfg.TokenMaxRequestsPerSecond)
	suite.Equal("6380", cfg.RedisPort)
	suite.Equal("testhost", cfg.RedisHost)
	suite.Equal("testpassword", cfg.RedisPassword)
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
