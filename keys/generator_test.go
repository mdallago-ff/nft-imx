package keys

import (
	"github.com/stretchr/testify/suite"
	"testing"
)

type UnitTestSuite struct {
	suite.Suite
}

func (s *UnitTestSuite) TestGenerator() {
	pair, err := CreateKeys()
	s.Assertions.Nil(err)
	s.Assertions.NotEmpty(pair.Private)
	s.Assertions.NotEmpty(pair.Public)
	s.Assertions.NotEmpty(pair.Address)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
