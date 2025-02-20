package api

import(
	"errors"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

var (
	ErrNotFound 		= errors.New("item not found")
	ErrUnauthorized 	= errors.New("not authorized")
	ErrServer		 	= errors.New("server identified error")
	ErrHTTPForbiden		= errors.New("forbiden request")
)

var childLogger = log.With().Str("go-core", "api").Logger()

type ApiService struct {
}

func NewRestApiService() *ApiService{
	childLogger.Debug().Msg("*** NewRestApiService")

	return &ApiService {
	}
}

func (a *ApiService) CallApi(ctx context.Context, 
							url string,
							method string,
							header_x_apigw_api_id *string,
							header_authorization *string,
							body interface{}) (interface{}, int, error){

	childLogger.Debug().Msg("CallApi")
	childLogger.Debug().Msg("--------------------------")
	childLogger.Debug().Interface("method:url:x_apigw_api_id : ", url + " " + method + " " + *header_x_apigw_api_id).Msg("")
	childLogger.Debug().Interface("body : ", body).Msg("")
	childLogger.Debug().Msg("--------------------------")

	transportHttp := &http.Transport{}
	
	client := http.Client{
		Transport: otelhttp.NewTransport(transportHttp),
		Timeout: time.Second * 29,
	}

	payload := new(bytes.Buffer) 
	if body != nil{
		json.NewEncoder(payload).Encode(body)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, payload)
	if err != nil {
		childLogger.Error().Err(err).Msg("error NewRequestWithContext")
		return nil, http.StatusInternalServerError ,errors.New(err.Error())
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	if (header_x_apigw_api_id != nil){
		req.Header.Add("x-apigw-api-id", *header_x_apigw_api_id)
	}
	if (header_authorization != nil){
		req.Header.Add("authorization", *header_authorization)
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		childLogger.Error().Err(err).Msg("error client.Do")
		return nil, http.StatusInternalServerError, errors.New(err.Error())
	}

	childLogger.Debug().Int("StatusCode :", resp.StatusCode).Msg("")
	
	switch (resp.StatusCode) {
		case 401:
			return nil, http.StatusUnauthorized, ErrUnauthorized
		case 403:
			return nil, http.StatusForbidden, ErrHTTPForbiden
		case 200:
		case 400:
			return nil, http.StatusNotFound, ErrNotFound
		case 404:
			return nil, http.StatusNotFound, ErrNotFound
		default:
			return nil, http.StatusInternalServerError, ErrServer
	}

	result := body
	err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
		childLogger.Error().Err(err).Msg("error no ErrUnmarshal")
		return nil, http.StatusInternalServerError, errors.New(err.Error())
    }

	return result, http.StatusOK ,nil
}