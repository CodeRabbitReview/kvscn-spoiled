package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

//Response combines response status code and body
type Response struct {
	Body       []byte
	StatusCode int
}

//API combines url and http.Client
type API struct {
	client *http.Client
	url    string
}

// NewAPI creates new instance or API by input url and
// http.Client, where http.Transport.MaxIdleConns and
// http.Transport.MaxConnsPerHost changed to 20000
// and http.Client.Timeout changed to 10 seconds
func NewAPI(url string) *API {
	var t = http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 20000
	t.MaxConnsPerHost = 20000
	return &API{
		client: &http.Client{Transport: t, Timeout: 10 * time.Second},
		url:    url,
	}
}

// GetAll sends request to the server by http.GET method
// if some errors appear function returns an error
// Method takes body and response and returns it in Response format
func (c *API) GetAll() (Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/api/", nil)
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return Response{
		Body:       data,
		StatusCode: resp.StatusCode,
	}, nil
}

// Delete sends delete request to server by http.MethodDelete
// Method takes param. It has to be a json key.
// It sends request with this key in body
// if some errors appear function returns an error
// Method takes body and response and returns it in Response format
func (c *API) Delete(param string) (Response, error) {
	req, err := http.NewRequest(http.MethodDelete, c.url+"/api/", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return Response{
		Body:       data,
		StatusCode: resp.StatusCode,
	}, nil
}

// GetByID sends get request to server by http.MethodGet
// Method takes param. It has to be a json key.
// It sends request with this key in body
// if some errors appear function returns an error
// Method takes body and response and returns it in Response format
func (c *API) GetByID(param string) (Response, error) {
	req, err := http.NewRequest(http.MethodGet, c.url+"/api/id", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return Response{
		Body:       data,
		StatusCode: resp.StatusCode,
	}, nil
}

// AddOrUpdate sends create request to server by http.MethodPost
// Method takes param. Format has to be like:{
//    "key":"key",
//    "entity": {
//		"entity": "entity"
//    }
//}.
// It sends request with this data in body
// if some errors appear function returns an error
// Method takes body and response and returns it in Response format
// Response.Body is always nil
func (c *API) AddOrUpdate(param string) (Response, error) {
	buf := bytes.NewBuffer([]byte(param))
	req, err := http.NewRequest(http.MethodPost, c.url+"/api/", buf)
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer func() {
		resp.Body.Close()
		req.Body.Close()
	}()
	return Response{
		Body:       nil,
		StatusCode: resp.StatusCode,
	}, nil
}
