package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

const BASEURL = "http://localhost:8080/"

func RequestNewKey() (string, error) {
	resp, err := resty.New().R().Get("http://localhost:8080/sessionkey")
	if err != nil {
		return "", err
	}

	var jsonResponse SessionKeyReponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	return jsonResponse.Key, nil
}

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
	key, err := RequestNewKey()
	assert.Nil(suite.T(), err)

	suite.ApiClient.SetTimeout(1 * time.Minute)
	suite.ApiClient.SetHostURL(BASEURL)
	suite.ApiClient.SetHeader("Accept", "application/json")
	suite.ApiClient.SetHeader("Content-Type", "application/json")
	suite.ApiClient.SetQueryParam("key", key)
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
	response, _ := suite.ApiClient.R().Get("pwd")
	var jsonResponse PWDResponse
	json.Unmarshal(response.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_PWDFailsWithNonExistingKey() {
	key := "scusasetelodicomaseimoltoitaliano"
	response, err := suite.ApiClient.R().SetQueryParam("key", key).Get("pwd")
	assert.Nil(suite.T(), err)
	var jsonResponse PWDResponse
	json.Unmarshal(response.Body(), &jsonResponse)

	assert.Equal(suite.T(), "error", jsonResponse.Outcome)
	assert.Equal(suite.T(), "", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDIntoOneSubdir() {
	CDResp, _ := suite.ApiClient.R().SetQueryParam("path", "test/").Get("cd")

	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get("pwd")
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDIntoManySubdir() {
	CDResp, _ := suite.ApiClient.R().SetQueryParam("path", "test/sub/user").Get("cd")
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get("pwd")
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test/sub/user", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDToRoot() {
	CDResp1, _ := suite.ApiClient.R().SetQueryParam("path", "test/sub/").Get("cd")
	var jsonCDResp1 OutcomeResponse
	json.Unmarshal(CDResp1.Body(), &jsonCDResp1)
	assert.Equal(suite.T(), "ok", jsonCDResp1.Outcome)

	CDResp, _ := suite.ApiClient.R().SetQueryParam("path", "/").Get("cd")
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get("pwd")
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_CDDotDot() {

	CDResp1, _ := suite.ApiClient.R().SetQueryParam("path", "test/sub/").Get("cd")
	var jsonCDResp1 OutcomeResponse
	json.Unmarshal(CDResp1.Body(), &jsonCDResp1)
	assert.Equal(suite.T(), "ok", jsonCDResp1.Outcome)

	CDResp, _ := suite.ApiClient.R().SetQueryParam("path", "..").Get("cd")
	var jsonCDResp OutcomeResponse
	json.Unmarshal(CDResp.Body(), &jsonCDResp)
	assert.Equal(suite.T(), "ok", jsonCDResp.Outcome)

	PWDResp, _ := suite.ApiClient.R().Get("pwd")
	var jsonResponse PWDResponse
	json.Unmarshal(PWDResp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
	assert.Equal(suite.T(), "/test", jsonResponse.Path)
}

func (suite *DataPuddleTestSuite) Test_MKDIRSuccess() {
	resp, _ := suite.ApiClient.R().SetQueryParam("path", "test/sub/").Get("mkdir")
	var jsonResponse OutcomeResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	_, err := os.Stat("storage/test/sub")
	assert.False(suite.T(), os.IsNotExist(err))

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
}

func (suite *DataPuddleTestSuite) Test_STORESuccessful() {
	file, err := ioutil.ReadFile("fixtures/testbig.json")
	assert.Nil(suite.T(), err)

	suite.ApiClient.R().SetQueryParam("path", "test/").Get("cd")

	resp, _ := suite.ApiClient.R().
		SetBody(file).
		SetQueryParam("filename", "big.json").
		Post("store")
	var jsonResponse OutcomeResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "ok", jsonResponse.Outcome)
}

func (suite *DataPuddleTestSuite) Test_STOREFailsIfFileAlreadyExists() {
	suite.ApiClient.R().SetQueryParam("path", "test/sub/").Get("cd")

	suite.ApiClient.R().
		SetBody(`{"username":"testuser"}`).
		SetQueryParam("filename", "fail.json").
		Post("store")

	resp, _ := suite.ApiClient.R().
		SetBody(`{"username":"testuser"}`).
		SetQueryParam("filename", "fail.json").
		Post("store")
	var jsonResponse OutcomeResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "error", jsonResponse.Outcome)
}

func (suite *DataPuddleTestSuite) Test_RETRIEVESuccess() {
	file, err := ioutil.ReadFile("fixtures/small.json")
	assert.Nil(suite.T(), err)

	suite.ApiClient.R().SetQueryParam("path", "test/").Get("cd")

	suite.ApiClient.R().
		SetBody(file).
		SetQueryParam("filename", "small.json").
		Post("store")

	resp, _ := suite.ApiClient.R().
		SetQueryParam("filename", "small.json").
		Get("retrieve")
	var jsonResponse RetrieveResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	assert.Equal(suite.T(), string(file), jsonResponse.File)
}

func (suite *DataPuddleTestSuite) Test_RETRIEVEFail() {
	resp, _ := suite.ApiClient.R().
		SetQueryParam("filename", "nonexistent.json").
		Get("retrieve")
	var jsonResponse RetrieveResponse
	json.Unmarshal(resp.Body(), &jsonResponse)

	assert.Equal(suite.T(), "error", jsonResponse.Outcome)
}
