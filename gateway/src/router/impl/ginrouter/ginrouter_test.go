package ginrouter_test

import (
	"encoding/json"
	"io"
	"main/logger"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"main/dataaccess/impl/mockdatastore"
	"main/router/impl/ginrouter"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *ginrouter.GinRouter

func init() {
	dbcfg := &mockdatastore.DBConfig{}
	mockDatastore, _ := mockdatastore.NewDatastore(dbcfg)
	router, _ = ginrouter.CreateRouter(mockDatastore)
}

func performRequest(r http.Handler, method, path string, payload io.Reader, username, pwd string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, payload)
	w := httptest.NewRecorder()
	if payload != nil {
		req.Header.Add("Content-Type", "application/json")
		//req.Header.Add("Content-Length", strconv.Itoa(len(payload.)))
	}
	if username != "" {
		req.SetBasicAuth(username, pwd)
	}

	r.ServeHTTP(w, req)
	return w
}

func Test0000Ping(t *testing.T) {
	// Build our expected body
	body := gin.H{
		"message": "pong",
	}

	w := performRequest(router.GetHandler(), "GET", "/ping", nil, "", "")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["message"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], value)
}

func Test00DeleteUnauthorized(t *testing.T) {
	w := performRequest(router.GetHandler(), "DELETE", "/services/v1/samurl/bicycleurl", nil, "", "sdfsd")
	// the request gives a 403
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func Test01Delete(t *testing.T) {
	w := performRequest(router.GetHandler(), "DELETE", "/services/v1/samurl/bicycleurl", nil, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	// Convert the JSON response to a map
	// var response map[string]string
	// err := json.Unmarshal([]byte(w.Body.String()), &response)
	// // Grab the value & whether or not it exists
	// value, exists := response["message"]
	// // Make some assertions on the correctness of the response.
	// assert.Nil(t, err)
	// assert.True(t, exists)
	// assert.Equal(t, body["message"], value)
}

func Test02PutTableUnexisting(t *testing.T) {
	// Build our expected body
	jsonResponse := "{\"message\":\"service do not exist for params: '','samurl','bicycleurl','','',0,0,0\"}"

	jsonFilePath := "../../examples/table.enabled.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "PUT", "/services/v1/samurl/bicycleurl", jsonFile, "samurl", "passw0rd")
	// the request gives a 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Equal(t, jsonResponse, w.Body.String())
}

func Test10Insert(t *testing.T) {
	// Build our expected body
	jsonFilePath := "../../examples/table.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl", jsonFile, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
}

func Test11InsertDuplicated(t *testing.T) {
	// Build our expected body
	jsonFilePath := "../../examples/table.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl", jsonFile, "samurl", "passw0rd")
	// the request gives a 400
	assert.Equal(t, http.StatusBadRequest, w.Code)
	// Convert the JSON response to a map
	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	// Grab the value & whether or not it exists
	value, exists := response["message"]
	// Make some assertions on the correctness of the response.
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.True(t, strings.Contains(value, "Duplicate"))
}

func Test12UnauthorizedInsert(t *testing.T) {
	// Build our expected body
	jsonFilePath := "../../examples/table.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "badusr", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl", jsonFile, "badusr", "passw0rd")
	// the request gives a 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

// curl -d "@tablecols.en.json" -X POST https://localhost:8443/services/v1/samurl/bicycleurl/colnames/en -ik -u samurl:ddd
func Test13ColnamesInsert(t *testing.T) {
	// Build our expected body
	csvResponse := "en,code,color,size,price,currency\n"

	jsonFilePath := "../../examples/tablecols.en.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl/colnames/en", jsonFile, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, csvResponse, w.Body.String())
}

func Test14ValuesInsert(t *testing.T) {
	// Build our expected body
	jsonResponse := "{\"count\":6}"

	jsonFilePath := "../../examples/tablevalues.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl/values", jsonFile, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, jsonResponse, w.Body.String())
}

func Test30ReadDraftCSV(t *testing.T) {
	// Build our expected body
	csvResponse := "'en','samurl','bicycleurl','bicycle models for summer 2020','summer,2020',5,0,1\n"

	w := performRequest(router.GetHandler(), "GET", "/services/v1/samurl/bicycleurl.csv", nil, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, csvResponse, w.Body.String())
}

func Test31ReadNotEnabled(t *testing.T) {
	body := gin.H{
		"message": "service do not exist for params: '','samurl','bicycleurl','','',0,0,2",
	}
	w := performRequest(router.GetHandler(), "GET", "/services/v1/samurl/bicycleurl", nil, "hkjh", "")
	// the request gives a 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["message"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], value)
}

func Test32ReadDraftColnamesCSV(t *testing.T) {
	// Build our expected body
	csvResponse := "en,code,color,size,price,currency\n"

	w := performRequest(router.GetHandler(), "GET", "/services/v1/samurl/bicycleurl/colnames/en.csv", nil, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, csvResponse, w.Body.String())
}

func Test33ReadDraftColnamesNotEnabled(t *testing.T) {
	body := gin.H{
		"message": "service do not exist for params: 'en','samurl','bicycleurl','','',0,0,2",
	}
	w := performRequest(router.GetHandler(), "GET", "/services/v1/samurl/bicycleurl/colnames/en.csv", nil, "hkjh", "")
	// the request gives a 400
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]string
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["message"]
	assert.Nil(t, err)
	assert.True(t, exists)
	assert.Equal(t, body["message"], value)
}

func Test40PutTableUnauthrized(t *testing.T) {
	// Build our expected body
	jsonResponse := "{\"status\":\"you are not allowed\"}"

	jsonFilePath := "../../examples/table.enabled.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "baduser", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "PUT", "/services/v1/samurl/bicycleurl", jsonFile, "baduser", "passw0rd")
	// the request gives a 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
	assert.Equal(t, jsonResponse, w.Body.String())
}

func Test41PutTable(t *testing.T) {
	// Build our expected body
	csvResponse := "'en','samurl','bicycleurl','bicycle models for summer 2020','summer,2020',5,6,2\n"

	jsonFilePath := "../../examples/table.enabled.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "passw0rd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "PUT", "/services/v1/samurl/bicycleurl", jsonFile, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, csvResponse, w.Body.String())
}

func Test50ReadCSV(t *testing.T) {
	// Build our expected body
	csvResponse := "'en','samurl','bicycleurl','bicycle models for summer 2020','summer,2020',5,6,2\n"

	w := performRequest(router.GetHandler(), "GET", "/services/v1/samurl/bicycleurl.csv", nil, "samurl", "passw0rd")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, csvResponse, w.Body.String())
}

func Test51ReadEnabledCSV(t *testing.T) {
	// Build our expected body
	csvResponse := "'en','samurl','bicycleurl','bicycle models for summer 2020','summer,2020',5,6,2\n"

	w := performRequest(router.GetHandler(), "GET", "/services/v1/samurl/bicycleurl.csv", nil, "", "")
	// the request gives a 200
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, csvResponse, w.Body.String())
}
