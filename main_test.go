package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type response struct {
	Status    string `json:"status"`
	ShortName string `json:"short"`
	URL       string `json:"url"`
}

type testingSuite struct {
	TestBed *testing.T
	Router  *gin.Engine
}

type tinyURLTesting interface {
	Test(string)
}

var tb tinyURLTesting

func TestShortURLService(t *testing.T) {
	_tb := testingSuite{
		TestBed: t,
		Router:  setupRouter(),
	}
	tb = _tb

	tb.Test("https://www.google.com")
	for i := 0; i < 10; i++ {
		originalURL := "https://" + randomString(30)
		tb.Test(originalURL)
	}
	tb.Test("https://www.google.com")
}

func (tb testingSuite) Test(originalURL string) {
	shortURL := tb.requestShortURL(originalURL)
	URL := tb.requestLongURL(shortURL)
	assert.Equal(tb.TestBed, originalURL, URL)
}

func (tb testingSuite) requestShortURL(originalURL string) string {
	gin.SetMode(gin.TestMode)
	postBody, _ := json.Marshal(map[string]string{
		"url": originalURL,
	})
	requestBody := bytes.NewBuffer(postBody)
	req, _ := http.NewRequest("POST", "/add", requestBody)
	req.Header.Add("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	tb.Router.ServeHTTP(resp, req)
	assert.Equal(tb.TestBed, 200, resp.Code)

	body, _ := ioutil.ReadAll(resp.Body)
	var respBody response
	json.Unmarshal(body, &respBody)

	return respBody.ShortName
}

func (tb testingSuite) requestLongURL(shortURL string) string {
	// get
	req, _ := http.NewRequest("GET", "/"+shortURL, nil)

	resp := httptest.NewRecorder()
	tb.Router.ServeHTTP(resp, req)
	assert.Equal(tb.TestBed, 200, resp.Code)
	body, _ := ioutil.ReadAll(resp.Body)
	var respBody response
	json.Unmarshal(body, &respBody)
	return respBody.URL
}
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}
