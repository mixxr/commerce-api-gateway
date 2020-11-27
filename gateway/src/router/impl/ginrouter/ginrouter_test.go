package ginrouter_test

import (
	"encoding/json"
	"io"
	"main/logger"
	"net/http"
	"net/http/httptest"
	"os"
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

func TestPing(t *testing.T) {
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

func Test00Insert(t *testing.T) {
	// Build our expected body
	jsonFilePath := "../../examples/table.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "samurl", "sdfsd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl", jsonFile, "samurl", "sdfsd")
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

func Test01BadInsert(t *testing.T) {
	// Build our expected body
	jsonFilePath := "../../examples/table.json"
	jsonFile, _ := os.Open(jsonFilePath)

	logger.AppLogger.Info("test", "test", "Successfully opened", jsonFilePath, "badusr", "sdfsd")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	w := performRequest(router.GetHandler(), "POST", "/services/v1/samurl/bicycleurl", jsonFile, "badusr", "sdfsd")
	// the request gives a 401
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}
