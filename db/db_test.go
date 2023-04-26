package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
	"nft/models"
	"nft/test"
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
	err := s.migrations.Reset(context.Background())
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

	newUser := test.CreateDummyUser(id, "test@test.com")

	err = s.db.CreateUser(newUser)
	s.Assertions.Nil(err)

	user, err = s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(user)
	s.Assertions.Equal(user, newUser)
}

func (s *UnitTestSuite) TestUpdateUser() {
	id := uuid.New()
	user, err := s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.Nil(user)

	newUser := test.CreateDummyUser(id, "test@test.com")
	err = s.db.CreateUser(newUser)
	s.Assertions.Nil(err)

	user, err = s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(user)
	s.Assertions.Equal(user.StarkKey, "")

	user.StarkKey = "random key"
	err = s.db.UpdateUser(user)
	s.Assertions.Nil(err)

	user, err = s.db.GetUser(id)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(user)
	s.Assertions.Equal(user.StarkKey, "random key")
}

func (s *UnitTestSuite) TestGetUserByMail() {
	mail := uuid.NewString() + "@test.com"
	user, err := s.db.GetUserByMail(mail)
	s.Assertions.Nil(err)
	s.Assertions.Nil(user)

	newUser := &models.User{
		ID:     uuid.New(),
		ApiKey: uuid.NewString(),
		Mail:   mail,
	}

	err = s.db.CreateUser(newUser)
	s.Assertions.Nil(err)

	user, err = s.db.GetUserByMail(mail)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(user)
	s.Assertions.Equal(user, newUser)
}

func (s *UnitTestSuite) TestCreateCollection() {
	id := uuid.New()
	collection, err := s.db.GetCollection(id)
	s.Assertions.Nil(err)
	s.Assertions.Nil(collection)

	userID := uuid.New()
	newCollection := test.CreateDummyCollection(id, userID, "test address")

	err = s.db.CreateCollection(newCollection)
	s.Assertions.Nil(err)

	collection, err = s.db.GetCollection(id)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(collection)
	s.Assertions.Equal(collection, newCollection)
}

func (s *UnitTestSuite) TestCreateToken() {
	id := uuid.New()
	token, err := s.db.GetToken(id)
	s.Assertions.Nil(err)
	s.Assertions.Nil(token)

	collectionID := uuid.New()
	newToken := test.CreateDummyToken(id, collectionID, "1")

	err = s.db.CreateToken(newToken)
	s.Assertions.Nil(err)

	token, err = s.db.GetToken(id)
	s.Assertions.Nil(err)
	s.Assertions.NotNil(token)
	s.Assertions.Equal(token, newToken)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
