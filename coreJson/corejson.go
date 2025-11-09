package coreJson

import (
	"io"
	"errors"
	"net/http"
	"encoding/json"

	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("component","go-core").Str("package", "coreJson").Logger()

type JSONResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type CoreJson struct{
}

// About convert JSON
func (c *CoreJson) ReadJSON(r *http.Request, w http.ResponseWriter, data interface{}) error {
	childLogger.Debug().Str("func","ReadJSON").Send()

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

// About convert JSON
func (c *CoreJson) WriteJSON(w http.ResponseWriter, code int, data interface{}) error {
	childLogger.Debug().Str("func","WriteJSON").Send()
	
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

type APIError struct {
	StatusCode	int  `json:"statusCode"`
	Msg			string `json:"message"`
	TraceId		string `json:"request-id,omitempty"`
}

func (e *APIError) Error() string {
	return e.Msg
}

func (e *APIError) NewAPIError(	err error, 
								traceId string, 
								status ...int) APIError {
	childLogger.Debug().Str("func","NewAPIError").Send()

	statusCode := http.StatusBadRequest
	
	if len(status) > 0 {
		statusCode = status[0]
	}

	return APIError{
		StatusCode: statusCode,
		Msg:		err.Error(),
		TraceId:	traceId,
	}
}