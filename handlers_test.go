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

func (suite *DataPuddleTestSuite) Test_CDIntoOneSubdir() {
	os.MkdirAll("storage/test", 0777)
	defer os.RemoveAll("storage/test")

	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", key, "test/"))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", key))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDIntoManySubdir() {
	os.MkdirAll("storage/test/sub/user", 0777)
	defer os.RemoveAll("storage/test")

	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", key, "test/sub/user"))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", key))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test/sub/user", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDToRoot() {
	os.MkdirAll("storage/test/sub", 0777)
	defer os.RemoveAll("storage/test")

	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	CDResp1, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", key, "test/sub/"))
	var jsonCDResp1 OutcomeResponse
	json.Unmarshal(CDResp1.Body(), &jsonCDResp1)
	assert.Equal(suite.T(), "ok", jsonCDResp1.Outcome)

	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", key, "/"))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", key))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDDotDot() {
	os.MkdirAll("storage/test/sub", 0777)
	defer os.RemoveAll("storage/test")

	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	CDResp1, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", key, "test/sub/"))
	var jsonCDResp1 OutcomeResponse
	json.Unmarshal(CDResp1.Body(), &jsonCDResp1)
	assert.Equal(suite.T(), "ok", jsonCDResp1.Outcome)

	CDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "cd", key, ".."))
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s", BASEURL, "pwd", key))
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_MKDIRSuccess() {
	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	resp, _ := suite.ApiClient.R().Get(fmt.Sprintf("%s/%s?key=%s&path=%s", BASEURL, "mkdir", key, "test/sub/"))
	var jsonResponse OutcomeResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	_,err = os.Stat("storage/test/sub")
	assert.False(suite.T(), os.IsNotExist(err))

	os.RemoveAll("storage/test")
	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
}
