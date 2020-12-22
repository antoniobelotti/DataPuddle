package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"os"
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
	apikey    string
}

func (suite *DataPuddleTestSuite) SetupTest() {
	suite.ApiClient = resty.New()
	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)
	suite.apikey = key
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(DataPuddleTestSuite))
}

func (suite *DataPuddleTestSuite) SetupSuite() {
	os.MkdirAll("storage/test/sub/user", 0777)
}

func (suite *DataPuddleTestSuite) TearDownSuite() {
	os.RemoveAll("storage/test")
}

func (suite *DataPuddleTestSuite) Test_IndexReturns_200_ok() {
	resp, err := suite.ApiClient.R().Get(BASEURL)
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), 200, resp.StatusCode())
}

func (suite *DataPuddleTestSuite) Test_PWDWorksWithExistingKey() {
	response, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", suite.apikey))
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

func (suite *DataPuddleTestSuite) Test_CDIntoOneSubdir() {
	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", suite.apikey, "test/"))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", suite.apikey))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDIntoManySubdir() {
	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", suite.apikey, "test/sub/user"))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", suite.apikey))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test/sub/user", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDToRoot() {
	CDResp1, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", suite.apikey, "test/sub/"))
	var jsonCDResp1 OutcomeResponse
	json.Unmarshal(CDResp1.Body(), &jsonCDResp1)
	assert.Equal(suite.T(), "ok", jsonCDResp1.Outcome)

	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", suite.apikey, "/"))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", suite.apikey))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDDotDot() {

	CDResp1, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", suite.apikey, "test/sub/"))
	var jsonCDResp1 OutcomeResponse
	json.Unmarshal(CDResp1.Body(), &jsonCDResp1)
	assert.Equal(suite.T(), "ok", jsonCDResp1.Outcome)

	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", suite.apikey, ".."))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", suite.apikey))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_MKDIRSuccess() {
	resp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "mkdir", suite.apikey, "test/sub/"))
	var jsonResponse OutcomeResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	_, err := os.Stat("storage/test/sub")
	assert.False(suite.T(), os.IsNotExist(err))

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
}
