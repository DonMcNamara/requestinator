package main

import (
	"crypto/rand"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"sync"

	"time"

	"encoding/base32"

	"strings"

	"gopkg.in/gin-gonic/gin.v1"
)

func main() {
	portNum := flag.Int("p", 8080, "port number to listen on")
	flag.Parse()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	context := context{bins: make(map[string]bin)}
	context.bins["test"] = bin{ID: "test"}
	router.POST("/api/v1/bins", context.createBin)
	router.GET("/api/v1/bins/:id", context.binDetails)
	router.GET("/api/v1/bins/:id/requests", context.binRequestDetails)
	router.Any("/bin/:id", context.requestHandler)

	address := ":" + strconv.Itoa(*portNum)
	log.Println("Listening on", address)
	router.Run()
}

type context struct {
	sync.Mutex // course-grained mutex for now
	bins       map[string]bin
}

type bin struct {
	ID           string    `json:"name"`
	RequestCount int       `json:"request_count"`
	Requests     []request `json:"-"`
}

type request struct {
	ContentLength int64       `json:"content_length"`
	ContentType   string      `json:"content_type"`
	Time          float64     `json:"time"`
	Method        string      `json:"method"`
	Body          string      `json:"body"`
	Header        http.Header `json:"headers"`
	QueryString   url.Values  `json:"query_string"`
	FormData      url.Values  `json:"form_data"`
}

func (binContext *context) createBin(c *gin.Context) {
	bin := bin{ID: createBinID()}
	binContext.bins[bin.ID] = bin
	c.JSON(http.StatusOK, bin)
}

func createBinID() string {
	bytes := make([]byte, 6)
	safeRandom(bytes)
	encoded := base32.StdEncoding.EncodeToString(bytes)
	return strings.ToLower(strings.Replace(encoded, "=", "", -1))
}

func safeRandom(dest []byte) {
	if _, err := rand.Read(dest); err != nil {
		panic(err)
	}
}

func (binContext *context) binDetails(c *gin.Context) {
	binID := c.Param("id")
	bin, ok := binContext.bins[binID]
	if ok {
		c.JSON(http.StatusOK, bin)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "bin not found"})
	}
}

func (binContext *context) binRequestDetails(c *gin.Context) {
	binID := c.Param("id")
	bin, ok := binContext.bins[binID]
	if ok {
		c.JSON(http.StatusOK, bin.Requests)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"message": "bin not found"})
	}
}

func (binContext *context) requestHandler(c *gin.Context) {
	binID := c.Param("id")
	bin, ok := binContext.bins[binID]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "bin not found"})
		return
	}

	req, err := fromContext(c)
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to parse request: "+err.Error())
		return
	}

	binContext.Lock()
	defer binContext.Unlock()
	bin.RequestCount++
	bin.Requests = append(bin.Requests, req)
	binContext.bins[binID] = bin
	c.String(http.StatusOK, "ok")
}

func fromContext(c *gin.Context) (request, error) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		return request{}, err
	}

	return request{
		Method:        c.Request.Method,
		ContentLength: c.Request.ContentLength,
		Header:        c.Request.Header,
		Time:          float64(time.Now().UnixNano()) / 1e9,
		ContentType:   c.ContentType(),
		FormData:      c.Request.Form,
		QueryString:   c.Request.URL.Query(),
		Body:          string(body),
	}, nil

}
