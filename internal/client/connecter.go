package client

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"time"
)

const address = "http://localhost:8080/api"

//Response combines response status code and body
type Response struct {
	Body       []byte
	StatusCode int
}

// GetAll sends request to the server by http.GET method
// if some errors appear function returns an error
// Method takes body and response and returns it in Response format
func GetAll() (Response, error) {
	resp, err := http.Get(address + "/")
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
func Delete(param string) (Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodDelete, address+"/", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return Response{}, err
	}
	resp, err := client.Do(req)
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
func GetByID(param string) (Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodGet, address+"/id", bytes.NewBuffer([]byte(param)))
	if err != nil {
		return Response{}, err
	}
	resp, err := client.Do(req)
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
func AddOrUpdate(param string) (Response, error) {
	var t = http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 20000
	t.MaxConnsPerHost = 20000
	defer t.CloseIdleConnections()
	var client = &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
	buf := bytes.NewBuffer([]byte(param))
	req, err := http.NewRequest(http.MethodPost, address+"/", buf)
	if err != nil {
		return Response{}, err
	}
	resp, err := client.Do(req)
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
