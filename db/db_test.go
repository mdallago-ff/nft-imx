package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"nft/models"
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
	err := s.migrations.Down(context.Background())
	s.Assertions.Nil(err)
}

func (s *UnitTestSuite) TestGetUser() {
	id := uuid.New()
	user, err := s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.Nil(user)
}

func (s *UnitTestSuite) TestCreateUser() {
	id := uuid.New()
	user, err := s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.Nil(user)

	newUser := &models.User{
		ID:     id,
		ApiKey: uuid.NewString(),
		Mail:   "test@test.com",
	}

	err = s.db.CreateUser(newUser)
	s.Assertions.Nil(err)

	user, err = s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(user)
	s.Assertions.Equal(user, newUser)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
