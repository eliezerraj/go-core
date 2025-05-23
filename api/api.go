package api

import(
	"io"
	"io/ioutil"
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
	Client		*http.Client
}

func NewRestApiService() *ApiService{
	childLogger.Debug().Str("func","NewRestApiService").Send()

	transportHttp := &http.Transport{
		MaxIdleConns:	100,
		MaxConnsPerHost: 100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:	10 * time.Second,
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(transportHttp),
		Timeout: 10 * time.Second,
	}

	return &ApiService {
		Client: &client,
	}
}

type HttpClient struct{
	Url 	string
	Method 	string
	Timeout time.Duration
	Headers *map[string]string
}

func (a *ApiService) CallRestApiV1(	ctx context.Context,
									client	*http.Client,		
									httpClient HttpClient,
									body interface{}) (interface{}, int, error){
	childLogger.Print("...................................")
	childLogger.Debug().Str("func","CallRestApiV1").
						Interface("client", client).
						Interface("httpClient", httpClient).
						Interface("body", body).Send()
	childLogger.Print("...................................")

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
		if key == "Host" || key == "host" {
			req.Host = value
		} else {
			req.Header.Set(key, value)
		}
	}

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, http.StatusServiceUnavailable, errors.New(err.Error())
	}

	defer func() {
		// Close body
		childLogger.Debug().Str("func", "CallRestApiV1").Msg("Body.Close !!!")
		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			childLogger.Error().Err(err).Send()
		}
		resp.Body.Close()
	}()

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