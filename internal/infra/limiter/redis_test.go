package limiter

import (
	"testing"
	"time"

	"github.com/alicebob/miniredis/v2"

	"github.com/stretchr/testify/suite"
)

type RedisTestSuite struct {
	suite.Suite
	redisServer *miniredis.Miniredis
	redis       *Redis
}

func (suite *RedisTestSuite) SetupTest() {
	var err error
	suite.redisServer, err = miniredis.Run()
	suite.Require().NoError(err)

	suite.redis = NewRedis(suite.redisServer.Addr(), "")
}

func (suite *RedisTestSuite) TearDownTest() {
	suite.redisServer.Close()
}

func (suite *RedisTestSuite) TestGet() {
	key := "test-key"
	value := "test-value"
	suite.redisServer.Set(key, value)

	result, err := suite.redis.Get(key)
	suite.NoError(err)
	suite.Equal(value, result)
}

func (suite *RedisTestSuite) TestGetNonExistentKey() {
	key := "non-existent-key"

	result, err := suite.redis.Get(key)
	suite.NoError(err)
	suite.Equal("", result)
}

func (suite *RedisTestSuite) TestSet() {
	key := "test-key"
	value := "test-value"
	expiration := 10 * time.Second

	err := suite.redis.Set(key, value, expiration)
	suite.NoError(err)

	result, err := suite.redis.Get(key)
	suite.NoError(err)
	suite.Equal(value, result)
}

func (suite *RedisTestSuite) TestIncr() {
	key := "test-key"

	err := suite.redis.Incr(key)
	suite.NoError(err)

	result, err := suite.redis.Get(key)
	suite.NoError(err)
	suite.Equal("1", result)
}

func (suite *RedisTestSuite) TestExpire() {
	key := "test-key"
	value := "test-value"
	expiration := 1 * time.Second

	err := suite.redis.Set(key, value, 0)
	suite.NoError(err)

	err = suite.redis.Expire(key, expiration)
	suite.NoError(err)

	time.Sleep(2 * time.Second)

	result, err := suite.redis.Get(key)
	suite.NoError(err)
	suite.Equal("test-value", result)
}

func TestRedisTestSuite(t *testing.T) {
	suite.Run(t, new(RedisTestSuite))
}
