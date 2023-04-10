package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"testing"
)

var dsn = "host=localhost user=postgres password=postgres dbname=nft_test port=5432 sslmode=disable"

type UnitTestSuite struct {
	suite.Suite
	db         *DB
	migrations *Migrations
}

func (s *UnitTestSuite) SetupTest() {
	migrations := NewMigrations(dsn)
	err := migrations.Up(context.Background())
	s.Assertions.Nil(err)
	db, err := NewDB(dsn)
	s.Assertions.Nil(err)
	s.db = db
	s.migrations = migrations
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	s.migrations.Down(context.Background())
}

func (s *UnitTestSuite) TestGetUser() {
	id := uuid.New()
	user, err := s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.Nil(user)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
