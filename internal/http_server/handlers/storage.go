package handlers

import (
	"github.com/mishaprokop4ik/storage/internal/models"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"log"
	"net/http"
	"strings"
)

type Storager interface {
	GetAll() (map[storage.Keyer]storage.Entitier, error)
	Get(key storage.Keyer) (storage.Entitier, error)
	Put(pair storage.Pair) error
	Delete(key storage.Keyer) error
}

type Storage struct {
	log     *log.Logger
	Storage Storager
}

func NewStorage(l *log.Logger, s Storager) *Storage {
	return &Storage{log: l, Storage: s}
}

func (s *Storage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.Replace(r.URL.String(), "/api/", "", 1)
	if strings.Count(url, "/") > 1 {
		SendResponse(w, Response{
			Data:       nil,
			StatusCode: http.StatusNotAcceptable,
		}, s.log)
		return
	}
	switch len([]rune(url)) == 0 {
	case true:
		if r.Method == http.MethodGet {
			s.GetAll(w, r)
		}

		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			s.Put(w, r)
		}
	case false:
		if r.Method == http.MethodDelete {
			s.Delete(w, r)
		}

		if r.Method == http.MethodGet {
			s.Get(w, r)
		}
	}
}

func (s *Storage) GetAll(w http.ResponseWriter, r *http.Request) {
	allData, err := s.Storage.GetAll()
	if err != nil {
		SendResponse(w, Response{
			Data:       err.Error(),
			StatusCode: http.StatusNotFound,
		}, s.log)
		return
	}
	var resp = "["
	var i = 0
	for _, d := range allData {
		v := d.JSON()[:len(d.JSON())]
		if i != len(allData)-1 {
			resp += string(v) + ", "
		} else {
			resp += string(v)
		}
		i++
	}
	resp += "]"

	SendResponse(w, Response{
		Data:       []byte(resp),
		StatusCode: http.StatusOK,
	}, s.log)
}

func (s *Storage) Get(w http.ResponseWriter, r *http.Request) {
	url := strings.Replace(r.URL.String(), "/api/", "", 1)
	param := models.NewKey(strings.Split(url, "/")[0])
	data, err := s.Storage.Get(param)
	if err != nil {
		SendResponse(w, Response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, s.log)
		return
	}

	SendResponse(w, Response{
		Data:       data.JSON(),
		StatusCode: http.StatusOK,
	}, s.log)
}

func (s *Storage) Put(w http.ResponseWriter, r *http.Request) {
	pair, err := GetBody(r)
	if err != nil {
		SendResponse(w, Response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, s.log)
	}

	err = s.Storage.Put(pair)
	if err != nil {
		SendResponse(w, Response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, s.log)
	}

	SendResponse(w, Response{
		Data:       nil,
		StatusCode: http.StatusOK,
	}, s.log)
}

func (s Storage) Delete(w http.ResponseWriter, r *http.Request) {
	url := strings.Replace(r.URL.String(), "/api/", "", 1)
	param := models.NewKey(strings.Split(url, "/")[0])
	err := s.Storage.Delete(param)
	if err != nil {
		SendResponse(w, Response{
			Data:       err.Error(),
			StatusCode: http.StatusInternalServerError,
		}, s.log)
		return
	}

	SendResponse(w, Response{
		StatusCode: http.StatusNoContent,
	}, s.log)
}
