package limiter

import (
	"testing"
	"time"

	"github.com/jonasjesusamerico/goexpert-rate-limiter/internal/infra/limiter/mock"
	"github.com/stretchr/testify/suite"
)

type LimiterTestSuite struct {
	suite.Suite
	limiter           *Limiter
	mockIPILimiter    *mock.ILimiter
	mockTokenILimiter *mock.ILimiter
}

// SetupTest inicializa os mocks e o objeto Limiter antes de cada teste
func (suite *LimiterTestSuite) SetupTest() {
	suite.mockIPILimiter = mock.NewILimiter(suite.T())
	suite.mockTokenILimiter = mock.NewILimiter(suite.T())
	suite.limiter = NewLimiter(5, 60, 10, 60, suite.mockIPILimiter, suite.mockTokenILimiter)
}

// TestAllowRequestByIP verifica se uma solicitação por IP é permitida
func (suite *LimiterTestSuite) TestAllowRequestByIP() {
	ip := "127.0.0.1"
	key := "ratelimit:ip:127.0.0.1"

	// Configurar expectativas para os métodos do mock
	suite.mockIPILimiter.EXPECT().Get(key).Return("4", nil)
	suite.mockIPILimiter.EXPECT().Incr(key).Return(nil)
	suite.mockIPILimiter.EXPECT().Expire(key, 60*time.Second).Return(nil)

	// Executar o teste
	result := suite.limiter.AllowRequest(ip, "")
	suite.True(result, "Expected request to be allowed")

	// Verificar se todas as expectativas foram atendidas
	suite.mockIPILimiter.AssertExpectations(suite.T())
}

// TestDenyRequestByIP verifica se uma solicitação por IP é negada
func (suite *LimiterTestSuite) TestDenyRequestByIP() {
	ip := "127.0.0.1"
	key := "ratelimit:ip:127.0.0.1"

	// Configurar expectativa de retorno que nega a solicitação
	suite.mockIPILimiter.EXPECT().Get(key).Return("5", nil)

	// Executar o teste
	result := suite.limiter.AllowRequest(ip, "")
	suite.False(result, "Expected request to be denied")

	// Verificar se todas as expectativas foram atendidas
	suite.mockIPILimiter.AssertExpectations(suite.T())
}

// TestAllowRequestByToken verifica se uma solicitação por token é permitida
func (suite *LimiterTestSuite) TestAllowRequestByToken() {
	token := "test-token"
	key := "ratelimit:token:test-token"

	// Configurar expectativas para os métodos do mock
	suite.mockTokenILimiter.EXPECT().Get(key).Return("9", nil)
	suite.mockTokenILimiter.EXPECT().Incr(key).Return(nil)
	suite.mockTokenILimiter.EXPECT().Expire(key, 60*time.Second).Return(nil)

	// Executar o teste
	result := suite.limiter.AllowRequest("", token)
	suite.True(result, "Expected request to be allowed")

	// Verificar se todas as expectativas foram atendidas
	suite.mockTokenILimiter.AssertExpectations(suite.T())
}

// TestDenyRequestByToken verifica se uma solicitação por token é negada
func (suite *LimiterTestSuite) TestDenyRequestByToken() {
	token := "test-token"
	key := "ratelimit:token:test-token"

	// Configurar expectativa de retorno que nega a solicitação
	suite.mockTokenILimiter.EXPECT().Get(key).Return("10", nil)

	// Executar o teste
	result := suite.limiter.AllowRequest("", token)
	suite.False(result, "Expected request to be denied")

	// Verificar se todas as expectativas foram atendidas
	suite.mockTokenILimiter.AssertExpectations(suite.T())
}

// TestLimiterTestSuite executa a suíte de testes
func TestLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
