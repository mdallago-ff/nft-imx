package server

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"nft/config"
	"nft/db"
	"nft/test"
	"testing"

	"github.com/go-chi/oauth"
	"github.com/google/uuid"
	"github.com/stretchr/testify/suite"
)

var dsn = "host=localhost user=postgres password=postgres dbname=nft_test port=5432 sslmode=disable"

func (s *UnitTestSuite) executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	s.server.Router.ServeHTTP(rr, req)
	return rr
}

func (s *UnitTestSuite) checkResponseCode(expected, actual int) {
	if expected != actual {
		s.T().Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

type UnitTestSuite struct {
	suite.Suite
	db         *db.DB
	migrations *db.Migrations
	server     *Server
}

func (s *UnitTestSuite) SetupTest() {
	migrations := db.NewMigrations(dsn)
	err := migrations.Up(context.Background())
	s.Assertions.Nil(err)
	newDB, err := db.NewDB(dsn)
	s.Assertions.Nil(err)
	s.db = newDB
	s.migrations = migrations
	settings := config.GetConfig()
	settings.DebugMode = true
	s.server = NewServer(settings, newDB, ImxDummy{})
	s.server.Configure()
}

func (s *UnitTestSuite) AfterTest(suiteName, testName string) {
	err := s.migrations.Reset(context.Background())
	s.Assertions.Nil(err)
}

func (s *UnitTestSuite) TestCreateUser() {
	var jsonStr = []byte(`{"mail":"test1@test.com"}`)
	req, _ := http.NewRequest("POST", "/users", bytes.NewBuffer(jsonStr))
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusCreated, response.Code)

	//"{\"id\":\"3e2a8125-7dd1-47b5-bea9-5a07c7132778\",\"email\":\"test1@test.com\",\"api_key\":\"7e418a84-4e15-4412-8165-fe31901e623a\",\"public\":\"0x04985d4379e537d6b1f9426477cd141bca57812b12bf6741320455f33f2eafe8db48f5672659fea3f183fc8c171a45b6783f51640e2edcc685805f36da12343156\",\"address\":\"0x6B138101C6fa0F30184B93585096d2F754782272\"}"
	objMap := map[string]string{}
	err := json.Unmarshal(response.Body.Bytes(), &objMap)
	s.Assertions.Nil(err)
	s.Assertions.NotEmpty(objMap["id"])
	s.Assertions.Equal("test1@test.com", objMap["email"])
	s.Assertions.NotEmpty(objMap["api_key"])
	s.Assertions.NotEmpty(objMap["public"])
	s.Assertions.NotEmpty(objMap["address"])
}

func (s *UnitTestSuite) TestCreateUserWithoutEmailShouldFail() {
	req, _ := http.NewRequest("POST", "/users", nil)
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusInternalServerError, response.Code)
}

func (s *UnitTestSuite) TestCreateCollection() {
	var jsonStr = []byte(`{"contract_address":"0x4958d0B91412eE2b8D715bF9279DCDB68e33d195", "collection_name":"prueba", "metadata_url":"https://gateway.pinata.cloud/ipfs/QmNj8NJwPbNGGv7HtjBii3TH1qa6yTmoJomvGth2rsXXyR", "fields": [ {"name":"name", "type": "text"} ]}`)
	req, _ := http.NewRequest("POST", "/collections", bytes.NewBuffer(jsonStr))
	ctx := req.Context()
	ctx = context.WithValue(ctx, oauth.CredentialContext, uuid.NewString())
	response := s.executeRequest(req.WithContext(ctx))
	s.checkResponseCode(http.StatusCreated, response.Code)

	objMap := map[string]string{}
	err := json.Unmarshal(response.Body.Bytes(), &objMap)
	s.Assertions.Nil(err)
	s.Assertions.Equal("0x4958d0B91412eE2b8D715bF9279DCDB68e33d195", objMap["contract_address"])
	s.Assertions.Equal("prueba", objMap["collection_name"])
	s.Assertions.NotEmpty(objMap["collection_id"])
}

func (s *UnitTestSuite) TestCreateToken() {
	user := test.CreateDummyUser(uuid.New(), "test")
	err := s.db.CreateUser(user)
	s.Assertions.Nil(err)
	collection := test.CreateDummyCollection(uuid.New(), user.ID, "address")
	err = s.db.CreateCollection(collection)
	s.Assertions.Nil(err)

	var jsonStr = []byte(`{"collection_id":"` + collection.ID.String() + `", "token_id": "1", "blueprint": "123456" }`)
	req, _ := http.NewRequest("POST", "/tokens", bytes.NewBuffer(jsonStr))

	ctx := req.Context()
	ctx = context.WithValue(ctx, oauth.CredentialContext, user.ID.String())
	response := s.executeRequest(req.WithContext(ctx))

	s.checkResponseCode(http.StatusCreated, response.Code)

	objMap := map[string]string{}
	err = json.Unmarshal(response.Body.Bytes(), &objMap)
	s.Assertions.Nil(err)
	s.Assertions.NotEmpty(objMap["token_id"])
}

func (s *UnitTestSuite) TestTransferToken() {
	user := test.CreateDummyUser(uuid.New(), "test")
	err := s.db.CreateUser(user)
	s.Assertions.Nil(err)
	collection := test.CreateDummyCollection(uuid.New(), user.ID, "address")
	err = s.db.CreateCollection(collection)
	s.Assertions.Nil(err)
	token := test.CreateDummyToken(uuid.New(), collection.ID, "1")
	err = s.db.CreateToken(token)
	s.Assertions.Nil(err)

	var jsonStr = []byte(`{"collection_id":"` + collection.ID.String() + `", "token_id":"` + token.ID.String() + `", "receiver_address": "0x18b1ceDC9803096D970f52260D1835F07D7e448C"}`)
	req, _ := http.NewRequest("POST", "/transfers", bytes.NewBuffer(jsonStr))

	ctx := req.Context()
	ctx = context.WithValue(ctx, oauth.CredentialContext, user.ID.String())
	response := s.executeRequest(req.WithContext(ctx))

	s.checkResponseCode(http.StatusCreated, response.Code)

	objMap := map[string]string{}
	err = json.Unmarshal(response.Body.Bytes(), &objMap)
	s.Assertions.Nil(err)
	s.Assertions.Equal(token.ID.String(), objMap["token_id"])
}

func (s *UnitTestSuite) TestCreateOrder() {
	user := test.CreateDummyUser(uuid.New(), "test")
	err := s.db.CreateUser(user)
	s.Assertions.Nil(err)
	collection := test.CreateDummyCollection(uuid.New(), user.ID, "address")
	err = s.db.CreateCollection(collection)
	s.Assertions.Nil(err)
	token := test.CreateDummyToken(uuid.New(), collection.ID, "1")
	err = s.db.CreateToken(token)
	s.Assertions.Nil(err)

	var jsonStr = []byte(`{"collection_id":"` + collection.ID.String() + `", "token_id":"` + token.ID.String() + `", "amount": "1000000"}`)
	req, _ := http.NewRequest("POST", "/orders", bytes.NewBuffer(jsonStr))

	ctx := req.Context()
	ctx = context.WithValue(ctx, oauth.CredentialContext, user.ID.String())
	response := s.executeRequest(req.WithContext(ctx))

	s.checkResponseCode(http.StatusCreated, response.Code)

	objMap := map[string]string{}
	err = json.Unmarshal(response.Body.Bytes(), &objMap)
	s.Assertions.Nil(err)
	s.Assertions.NotEmpty(objMap["order_id"])
}

func (s *UnitTestSuite) TestCreateDeposit() {
	user := test.CreateDummyUser(uuid.New(), "test")
	err := s.db.CreateUser(user)
	s.Assertions.Nil(err)

	var jsonStr = []byte(`{"amount_wei":"1000000000"}`)
	req, _ := http.NewRequest("POST", "/deposits", bytes.NewBuffer(jsonStr))

	ctx := req.Context()
	ctx = context.WithValue(ctx, oauth.CredentialContext, user.ID.String())
	response := s.executeRequest(req.WithContext(ctx))

	s.checkResponseCode(http.StatusCreated, response.Code)

	objMap := map[string]string{}
	err = json.Unmarshal(response.Body.Bytes(), &objMap)
	s.Assertions.Nil(err)
	s.Assertions.NotEmpty(objMap["tx_hash"])
}

func (s *UnitTestSuite) TestCreateCollectionWithoutParamsShouldFail() {
	req, _ := http.NewRequest("POST", "/collections", nil)
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusInternalServerError, response.Code)
}

func (s *UnitTestSuite) TestCreateTokenWithoutParamsShouldFail() {
	req, _ := http.NewRequest("POST", "/tokens", nil)
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusInternalServerError, response.Code)
}

func (s *UnitTestSuite) TestTransferTokenWithoutParamsShouldFail() {
	req, _ := http.NewRequest("POST", "/transfers", nil)
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusInternalServerError, response.Code)
}

func (s *UnitTestSuite) TestCreateOrderWithoutParamsShouldFail() {
	req, _ := http.NewRequest("POST", "/orders", nil)
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusInternalServerError, response.Code)
}

func (s *UnitTestSuite) TestCreateDepositWithoutParamsShouldFail() {
	req, _ := http.NewRequest("POST", "/deposits", nil)
	response := s.executeRequest(req)
	s.checkResponseCode(http.StatusInternalServerError, response.Code)
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}
