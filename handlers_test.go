package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

const BASEURL = "http://localhost:8080/"

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
	resp, err := suite.ApiClient.R().Get(BASEURL)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode())
}

func (suite *DataPuddleTestSuite) Test_SessionKeyGenerationIsSuccessful() {
	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 32, len(key))
}

func (suite *DataPuddleTestSuite) Test_PWDWorksWithExistingKey() {
	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	response, err := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", key))
	assert.Nil(suite.T(), err)
	var jsonResponse PWDResponse
	json.Unmarshal(response.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_PWDFailsWithNonExistingKey() {
	key := "scusasetelodicomaseimoltoitaliano"
	response, err := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", key))
	assert.Nil(suite.T(), err)
	var jsonResponse PWDResponse
	json.Unmarshal(response.Body(), &jsonResponse)

	assert.Equal(suite.T(), "error", jsonResponse.Outcome)
	assert.Equal(suite.T(), "", jsonResponse.Path)
}

func RequestNewKey() (string, error) {
	resp, err := resty.New().R().Get("http://localhost:8080/sessionkey")
	if err != nil {
		return "", err
	}

	var jsonResponse SessionKeyReponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	return jsonResponse.Key, nil
}
