package handlers

import (
	"encoding/json"
	"fmt"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"html/template"
	"net/http"
	"sort"
	"strings"
)

var indexPath = "./internal/http_server/handlers/static/index.gohtml"

// Storager should implement 4 methods of common operations:
// getting all data, getting by id, putting new data, deleting
// Storager should implement:
// GetAll that returns all data from Storage
// value has map[storage.Keyer]storage.Entitier
// if some errors appears - return error
// Get parameter is a storage.Keyer
// returns storage.Entitier and error if appears
// Put insert new data or update old
// if some errors appears return it
// Delete parameter is a storage.Keyer
// if some errors appears return it
type Storager interface {
	GetAll() (map[storage.Keyer]storage.Entitier, error)
	Get(key storage.Keyer) (storage.Entitier, error)
	Put(pair storage.Pair) error
	Delete(key storage.Keyer) error
}

// Storage is compacted 2 value
// log is a log.Logger
// and storage is a Storager
type Storage struct {
	storage Storager
}

// NewStorage is a constructor of Storage
func NewStorage(s Storager) *Storage {
	return &Storage{storage: s}
}

// ServeHTTP should write reply headers and data to the ResponseWriter
// and then return. Returning signals that the request is finished; it
// is not valid to use the ResponseWriter or read from the
// Request.Body after or concurrently with the completion of the
// ServeHTTP call.
// ServeHTTP routes requested URL and call specific method of Storage
// URL prefix must be /api
// if URL is / and request method is http.MethodGet ServeHTTP calls GetAll method
// if URL is / and request method is http.MethodPut or http.MethodPost ServeHTTP calls GetAll method
// if URL is /:id and request method is http.MethodGet calls Get method
// id can be any value
// if URL is /:id and request method is http.MethodDelete calls Delete method
// id can be any value
// if URL is incorrect returns http.StatusNotFound and nothing in body
func (s *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zlog.Log.WithName("http server").Info("user call", "search host",
		r.Host, "user url", r.URL.User, "host", r.URL.Host)
	url := r.URL.String()

	if !strings.HasPrefix(url, "/api/") {
		sendResponse(w, response{
			Data:       nil,
			StatusCode: http.StatusNotFound,
		})
		return
	}

	if strings.Count(url, "/") > 2 {
		sendResponse(w, response{
			Data:       nil,
			StatusCode: http.StatusNotAcceptable,
		})
		return
	}

	url = strings.Replace(r.URL.String(), "/api", "", 1)

	if url == "/id" && r.Method == http.MethodGet {
		s.Get(w, r)
		return
	}

	if url == "/" && r.Method == http.MethodGet {
		s.GetAll(w, r)
		return
	}

	if url == "/out" && r.Method == http.MethodGet {
		s.OutHTML(w, r)
		return
	}

	if url == "/" && r.Method == http.MethodPost || r.Method == http.MethodPut {
		s.Put(w, r)
		return
	}

	if url == "/" && r.Method == http.MethodDelete {
		s.Delete(w, r)
		return
	}

	sendResponse(w, response{
		Data:       fmt.Errorf("not found action by input url"),
		StatusCode: http.StatusNotFound,
	})
}

// GetAll sends data to http.ResponseWriter in JSON format
// response is an array of JSON objects
// if no value in storage returns no data in storage error
func (s *Storage) GetAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	storageData, err := s.storage.GetAll()
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusNotFound,
		})
		return
	}

	respData := make([]json.RawMessage, len(storageData))
	var i int
	for _, v := range storageData {
		data := v.JSON()
		respData[i] = data
		i++
	}

	sort.Slice(respData, func(i, j int) bool {
		return string(respData[i]) < string(respData[j])
	})

	resp, err := json.Marshal(respData)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
	}

	sendResponse(w, response{
		Data:       resp,
		StatusCode: http.StatusOK,
	})
}

// Get sends data to http.ResponseWriter in JSON format
// response is an JSON object
// Method takes id from URL
// id can be any data
// if no value in storage returns no data in storage error
// Get takes first param from URL from http.Request
func (s *Storage) Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	pair, err := getPairFromBody(r)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	data, err := s.storage.Get(pair.Key)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusNotFound,
		})
		return
	}

	sendResponse(w, response{
		Data:       data.JSON(),
		StatusCode: http.StatusOK,
	})
}

// Put method takes data from http.Request body.
// If some error appears in getting data from http.Request returns http.StatusInternalServerError
// and error in JSON format.
// If some error appears in getting from Storager.Put returns http.StatusInternalServerError
// and error in JSON format.
// If everything is OK and input data stored returns http.StatusCreated
// and nothing in body.
// Where key value is string and if it has spaces between words or before or after
// Spaces before and after will be removed
// spaces between words will be changed to _ symbol
func (s *Storage) Put(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	pair, err := getPairFromBody(r)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	err = s.storage.Put(pair)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusBadRequest,
		})
		return
	}

	sendResponse(w, response{
		Data:       nil,
		StatusCode: http.StatusCreated,
	})
}

// Delete sends data to http.ResponseWriter in JSON format
// Method takes id from URL
// id can be any data
// it calls Storager.Delete
// if some error appears http.StatusInternalServerError
// if everything is OK returns http.StatusNoContent and nothing in body
func (s *Storage) Delete(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	pair, err := getPairFromBody(r)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
	err = s.storage.Delete(pair.Key)
	if err != nil {
		sendResponse(w, response{
			Data:       err.Error(),
			StatusCode: http.StatusNotFound,
		})
		return
	}

	sendResponse(w, response{
		StatusCode: http.StatusNoContent,
	})
}

// OutHTML takes all data from Storage
// takes html file and parse in this html file all got data.
// If no data in storage - will be printed such text.
// If there is any data in storage - will be printed table with data.
func (s *Storage) OutHTML(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	w.Header().Set("Content-Type", "text/html;charset=utf-8")
	t, err := template.ParseFiles(indexPath)
	if err != nil {
		sendResponse(w, response{
			Data:       nil,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}

	type data struct {
		Data map[string]interface{}
	}
	allStorageData, _ := s.storage.GetAll()
	var resp = data{
		Data: make(map[string]interface{}),
	}
	for k, e := range allStorageData {
		keyString, _ := k.Entity().(string)
		resp.Data[keyString] = e.Entity()
	}

	err = t.Execute(w, resp)
	if err != nil {
		sendResponse(w, response{
			Data:       nil,
			StatusCode: http.StatusInternalServerError,
		})
		return
	}
}
