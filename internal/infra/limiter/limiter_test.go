package limiter

import (
	"testing"
	"time"

	"github.com/ramonamorim/go-rate-limiter/internal/infra/limiter/mock"
	"github.com/stretchr/testify/suite"
)

type LimiterTestSuite struct {
	suite.Suite
	limiter           *Limiter
	mockIPILimiter    *mock.ILimiter
	mockTokenILimiter *mock.ILimiter
}

func (suite *LimiterTestSuite) SetupTest() {
	suite.mockIPILimiter = mock.NewILimiter(suite.T())
	suite.mockTokenILimiter = mock.NewILimiter(suite.T())
	suite.limiter = NewLimiter(5, 60, 10, 60, suite.mockIPILimiter, suite.mockTokenILimiter)
}

func (suite *LimiterTestSuite) TestAllowRequestByIP() {
	ip := "127.0.0.1"
	key := "ratelimit:ip:127.0.0.1"

	suite.mockIPILimiter.EXPECT().Get(key).Return("4", nil)
	suite.mockIPILimiter.EXPECT().Incr(key).Return(nil)
	suite.mockIPILimiter.EXPECT().Expire(key, 60*time.Second).Return(nil)

	result := suite.limiter.AllowRequest(ip, "")
	suite.True(result, "Expected request to be allowed")

	suite.mockIPILimiter.AssertExpectations(suite.T())
}

func (suite *LimiterTestSuite) TestDenyRequestByIP() {
	ip := "127.0.0.1"
	key := "ratelimit:ip:127.0.0.1"

	suite.mockIPILimiter.EXPECT().Get(key).Return("5", nil)

	result := suite.limiter.AllowRequest(ip, "")
	suite.False(result, "Expected request to be denied")

	suite.mockIPILimiter.AssertExpectations(suite.T())
}

func (suite *LimiterTestSuite) TestAllowRequestByToken() {
	token := "test-token"
	key := "ratelimit:token:test-token"

	suite.mockTokenILimiter.EXPECT().Get(key).Return("9", nil)
	suite.mockTokenILimiter.EXPECT().Incr(key).Return(nil)
	suite.mockTokenILimiter.EXPECT().Expire(key, 60*time.Second).Return(nil)

	result := suite.limiter.AllowRequest("", token)
	suite.True(result, "Expected request to be allowed")

	suite.mockTokenILimiter.AssertExpectations(suite.T())
}

func (suite *LimiterTestSuite) TestDenyRequestByToken() {
	token := "test-token"
	key := "ratelimit:token:test-token"

	suite.mockTokenILimiter.EXPECT().Get(key).Return("10", nil)

	result := suite.limiter.AllowRequest("", token)
	suite.False(result, "Expected request to be denied")

	suite.mockTokenILimiter.AssertExpectations(suite.T())
}

func TestLimiterTestSuite(t *testing.T) {
	suite.Run(t, new(LimiterTestSuite))
}
