package client

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

//Response combines response status code and body
type Response struct {
	Body       []byte
	StatusCode int
}

//API combines server url
//and http.Client
type API struct {
	client *http.Client
	url    string
}

// NewAPI creates new instance or API by input url and
// http.Client, where http.Transport.MaxIdleConns and
// http.Transport.MaxConnsPerHost changed to 20000
// and http.Client.Timeout changed to 10 seconds
// http.Client will send data by https protocol
// it no files by certPath NewAPI will send data by http
// instead of https
func NewAPI(url, certPath string) *API {
	var t = http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 20000
	t.MaxConnsPerHost = 20000
	_, err := os.Stat(certPath)
	if errors.Is(err, os.ErrNotExist) {
		zlog.Log.WithName("connector").
			Info("create client with http only")
		return &API{
			client: &http.Client{Transport: t, Timeout: 10 * time.Second},
			url:    url,
		}
	}
	caCert, err := ioutil.ReadFile(certPath)
	if err != nil {
		zlog.Log.WithName("connector").
			Error(err, "can not read certificate")
	}
	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)
	t.TLSClientConfig = &tls.Config{
		MinVersion: tls.VersionTLS12,
		RootCAs:    caCertPool,
	}
	return &API{
		client: &http.Client{Transport: t, Timeout: 10 * time.Second},
		url:    url,
	}
}

// GetAll sends request to the server by http.GET method
// if some errors appear function returns an error
// Method takes body and response and returns it in Response format
func (c *API) GetAll() (Response, error) {
	zlog.Log.WithName("http client").
		Info("get all data from server", "url", c.url)
	req, err := http.NewRequest(http.MethodGet, c.url+"/api/", nil)
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		zlog.Log.WithName("http client").
			Error(err, "can not send request to server")
		return Response{}, err
	}

	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	if resp.StatusCode != http.StatusOK {
		return Response{
				Body:       nil,
				StatusCode: resp.StatusCode,
			}, fmt.Errorf("incorrect status code want: %d; get: %d",
				http.StatusOK, resp.StatusCode)
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
	zlog.Log.WithName("http client").
		Info("delete data from server", "url", c.url,
			"id", param)
	req, err := http.NewRequest(http.MethodDelete, c.url+"/api/", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		zlog.Log.WithName("http client").
			Error(err, "can not send request to server")
		return Response{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode != http.StatusNoContent {
		return Response{
				Body:       nil,
				StatusCode: resp.StatusCode,
			}, fmt.Errorf("incorrect status code want: %d; get: %d",
				http.StatusNoContent, resp.StatusCode)
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
	zlog.Log.WithName("http client").
		Info("get data from server", "url", c.url,
			"id", param)
	req, err := http.NewRequest(http.MethodGet, c.url+"/api/id", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		zlog.Log.WithName("http client").
			Error(err, "can not send request to server")
		return Response{}, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}
	if resp.StatusCode != http.StatusOK {
		return Response{
				Body:       nil,
				StatusCode: resp.StatusCode,
			}, fmt.Errorf("incorrect status code want: %d; get: %d",
				http.StatusOK, resp.StatusCode)
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
	zlog.Log.WithName("http client").
		Info("create or update data", "url", c.url)
	buf := bytes.NewBuffer([]byte(param))
	req, err := http.NewRequest(http.MethodPost, c.url+"/api/", buf)
	if err != nil {
		return Response{}, err
	}
	resp, err := c.client.Do(req)
	if err != nil {
		zlog.Log.WithName("http client").
			Error(err, "can not send request to server")
		return Response{}, err
	}
	if resp.StatusCode != http.StatusCreated {
		return Response{
				Body:       nil,
				StatusCode: resp.StatusCode,
			}, fmt.Errorf("incorrect status code want: %d; get: %d",
				http.StatusCreated, resp.StatusCode)
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
