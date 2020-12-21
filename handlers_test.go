package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

func init() {
	app := App{router: SetUpRouter()}
	go app.Run(fmt.Sprintf(":%s", "8080"))
}

type DataPuddleTestSuite struct {
	suite.Suite
	ApiClient *resty.Client
}

func (suite *DataPuddleTestSuite) SetupTest() {
	suite.ApiClient = resty.New()
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DataPuddleTestSuite))
}

func (suite *DataPuddleTestSuite) Test_IndexReturns_200_ok() {
	resp, err := suite.ApiClient.R().Get("http://localhost:8080/")
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode())
}

func (suite *DataPuddleTestSuite) Test_SessionKeyGenerationIsSuccessful() {
	resp, err := suite.ApiClient.R().Get("http://localhost:8080/sessionkey")
	assert.Nil(suite.T(), err)

	var jsonResponse SessionKeyReponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), 32, len(jsonResponse.Key))
}
