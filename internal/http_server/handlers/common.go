package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/mishaprokop4ik/storage/internal/models"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"io/ioutil"
	"log"
	"net/http"
)

type Response struct {
	Data       interface{} `json:"response"`
	StatusCode int         `json:"-"`
}

func SendResponse(w http.ResponseWriter, data Response, logger *log.Logger) {
	if v, ok := data.Data.([]byte); ok {
		if data.Data != nil {
			if _, err := w.Write(v); err != nil {
				logger.Fatal(err)
			}
		}
		w.WriteHeader(data.StatusCode)
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		if _, err = w.Write([]byte(err.Error())); err != nil {
			logger.Fatal(err)
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if data.Data != nil {
		if _, err = w.Write(resp); err != nil {
			logger.Fatal(err)
		}
	}
	w.WriteHeader(data.StatusCode)
}

func GetBody(r *http.Request) (storage.Pair, error) {
	defer r.Body.Close()
	type pair struct {
		Key    interface{}
		Entity interface{}
	}
	var p pair
	bodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return storage.Pair{}, err
	}
	r.Body.Close()
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return storage.Pair{}, err
	}
	return storage.Pair{
		Key:    models.NewKey(p.Key),
		Entity: models.NewEntity(p.Entity, bodyBytes),
	}, nil
}
