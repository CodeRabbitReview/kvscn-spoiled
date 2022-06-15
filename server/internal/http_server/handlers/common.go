package handlers

import (
	"bytes"
	"encoding/json"
	zlog "github.com/mishaprokop4ik/storage/internal/log"
	models2 "github.com/mishaprokop4ik/storage/internal/models"
	"github.com/mishaprokop4ik/storage/internal/storage"
	"io/ioutil"
	"net/http"
)

type response struct {
	Data       interface{} `json:"response"`
	StatusCode int         `json:"-"`
}

func sendResponse(w http.ResponseWriter, data response) {
	w.WriteHeader(data.StatusCode)
	if v, ok := data.Data.([]byte); ok {
		if data.Data != nil {
			if _, err := w.Write(v); err != nil {
				zlog.Log.WithName("http server").
					Error(err, "can not send response body")
				return
			}
		}
		return
	}

	resp, err := json.Marshal(data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err = w.Write([]byte(err.Error())); err != nil {
			zlog.Log.WithName("http server").
				Error(err, "can not send response header")
			return
		}
		return
	}

	if data.Data != nil {
		if _, err = w.Write(resp); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			zlog.Log.WithName("http server").
				Error(err, "can not send response body")
			return
		}
	}
}

func getPairFromBody(r *http.Request) (storage.Pair, error) {
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
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		return storage.Pair{}, err
	}
	e, err := models2.NewClearEntity(p.Entity, bodyBytes)
	if err != nil {
		return storage.Pair{}, err
	}
	return storage.Pair{
		Key:    models2.NewKey(p.Key),
		Entity: e,
	}, nil
}
