package coreJson

import (
	"io"
	"errors"
	"net/http"
	"encoding/json"
)

type CoreJson struct{
}

type JSONResponse struct {
	Error   bool        `json:"error"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (c *CoreJson) ReadJSON(r *http.Request, w http.ResponseWriter, data interface{}) error {
	maxBytes := 1024 * 1024 // 1 Mb

	r.Body = http.MaxBytesReader(w, r.Body, int64(maxBytes))

	decoded := json.NewDecoder(r.Body)
	err := decoded.Decode(data)
	if err != nil {
		return err
	}

	err = decoded.Decode(&struct{}{})
	if err != io.EOF {
		return errors.New("body must contain only one JSON value")
	}

	return nil
}

func (c *CoreJson) WriteJSON(w http.ResponseWriter, code int, data interface{}) error {
	data_json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")
	
	_, err = w.Write(data_json)
	if err != nil {
		return err
	}
	return nil
}