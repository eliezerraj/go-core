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

var childLogger = log.With().Str("component","go-core").Str("package", "api").Logger()

var (
	ErrNotFound 		= errors.New("item not found")
	ErrUnauthorized 	= errors.New("not authorized")
	ErrServer		 	= errors.New("server identified error")
	ErrHTTPForbiden		= errors.New("forbiden request")
)

type ApiService struct {
}

func NewRestApiService() *ApiService{
	childLogger.Debug().Str("func","NewRestApiService").Send()

	return &ApiService {
	}
}

type HttpClient struct{
	Url 	string
	Method 	string
	Timeout time.Duration
	Headers *map[string]string
}

func (a *ApiService) CallApi(ctx context.Context, 
							url string,
							method string,
							header_x_apigw_api_id *string,
							header_authorization *string,
							header_x_resquest_id *string,
							body interface{}) (interface{}, int, error){
	childLogger.Print("...................................")
	childLogger.Debug().Str("func","CallApi").
						Interface("parameters", method + ":" +url + ":" + *header_x_apigw_api_id  + ":" + *header_x_resquest_id).
						Interface("body", body).Send()
	childLogger.Print("...................................")

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
		childLogger.Error().Err(err).Send()
		return nil, http.StatusBadGateway, errors.New(err.Error())
	}

	req.Header.Add("Content-Type", "application/json;charset=UTF-8");
	if (header_x_apigw_api_id != nil){
		req.Header.Add("x-apigw-api-id", *header_x_apigw_api_id)
	}
	if (header_authorization != nil){
		req.Header.Add("authorization", *header_authorization)
	}
	if (header_x_resquest_id != nil){
		req.Header.Add("X-Request-Id", *header_x_resquest_id)
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, http.StatusServiceUnavailable, errors.New(err.Error())
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
		childLogger.Error().Err(err).Send()
		return nil, http.StatusInternalServerError, errors.New(err.Error())
    }

	return result, http.StatusOK ,nil
}

// Above call rest api 
func (a *ApiService) CallRestApi(ctx context.Context, 
							httpClient HttpClient,
							body interface{}) (interface{}, int, error){
	childLogger.Print("...................................")
	childLogger.Debug().Str("func","CallRestApi").
						Interface("httpClient", httpClient).
						Interface("body", body).Send()
	childLogger.Print("...................................")

	transportHttp := &http.Transport{}
	
	client := http.Client{
		Transport: otelhttp.NewTransport(transportHttp),
		Timeout: time.Second * httpClient.Timeout,
	}

	payload := new(bytes.Buffer) 
	if body != nil{
		json.NewEncoder(payload).Encode(body)
	}

	req, err := http.NewRequestWithContext(ctx, httpClient.Method, httpClient.Url, payload)
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, http.StatusBadGateway, errors.New(err.Error())
	}

	for key, value := range *httpClient.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, http.StatusServiceUnavailable, errors.New(err.Error())
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
		childLogger.Error().Err(err).Send()
		return nil, http.StatusInternalServerError, errors.New(err.Error())
    }

	return result, http.StatusOK ,nil
}