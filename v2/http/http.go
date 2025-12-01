package http

import(
	"io"
	"fmt"
	"io/ioutil"
	"errors"
	"net/http"
	"time"
	"encoding/json"
	"bytes"
	"syscall"
	"context"

	"github.com/rs/zerolog"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

type HttpService struct {
	client		*http.Client
	logger 		*zerolog.Logger
}

type HttpClientParameter struct{
	Url 	string
	Method 	string
	Timeout time.Duration
	Headers *map[string]string
	Body 	interface{}
}

// Above create a http service
func NewHttpService(appLogger *zerolog.Logger) *HttpService {
	logger := appLogger.With().
						Str("component", "go-core.v2.http").
						Logger()
	logger.Debug().
			Str("func","NewHttpService").Send()

	transportHttp := &http.Transport{
		MaxIdleConns:	100,
		MaxConnsPerHost: 100,
		MaxIdleConnsPerHost: 100,
		IdleConnTimeout:	100 * time.Second,
	}

	client := http.Client{
		Transport: otelhttp.NewTransport(transportHttp),
		Timeout: 10 * time.Second,
	}

	return &HttpService {
		client: &client,
		logger: &logger,
	}
}

// Above do a call via http service
func (d *HttpService) DoHttp(ctx context.Context,
							 httpClientParameter HttpClientParameter) (	interface{}, 
																		int, 
																		error){
	d.logger.Debug().
			 Ctx(ctx).
			 Str("func","DoHttp").Send()

	d.logger.Debug().
			Ctx(ctx).
			Interface("Body", httpClientParameter.Body).Send()

	payload := new(bytes.Buffer) 
	if httpClientParameter.Body != nil{
		json.NewEncoder(payload).Encode(httpClientParameter.Body)
	}

	req, err := http.NewRequestWithContext(	ctx, 
											httpClientParameter.Method, 
											httpClientParameter.Url, 
											payload)
	if err != nil {
		d.logger.Error().
				 Ctx(ctx).
				 Err(err).Send()
		return nil, http.StatusInternalServerError, errors.New(err.Error())
	}

	for key, value := range *httpClientParameter.Headers {
		if key == "Host" || key == "host" {
			req.Host = value
		} else {
			req.Header.Set(key, value)
		}
	}

	resp, err := d.client.Do(req.WithContext(ctx))
	if err != nil {
		if errors.Is(err, syscall.EPIPE) {
			d.logger.Error().
			         Ctx(ctx).
					 Err(err).
					 Msg("WARNING: BROKEN PIPE ERROR")
        } else if errors.Is(err, syscall.ECONNRESET)  {
			d.logger.Error().
					 Ctx(ctx).
					 Err(err).
					 Msg("CONNECTION RESET BY PIER")
		} else {
			d.logger.Error().
			         Ctx(ctx). 
					 Err(err).Send()
		}
		return nil, http.StatusInternalServerError, errors.New(err.Error())
	}

	defer func() {
		// Close body
		d.logger.Debug().
		         Ctx(ctx).
				 Msg("Body.Close !!!")

		if _, err := io.Copy(ioutil.Discard, resp.Body); err != nil {
			d.logger.Error().
					 Ctx(ctx).
					 Err(err).Send()
		}
		resp.Body.Close()
	}()

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
    if err != nil {
		d.logger.Error().
		         Ctx(ctx).
				 Err(err).Send()
		return nil, http.StatusInternalServerError, errors.New(err.Error())
    }

	// Any error message come from a go-core middleware has a field message, according to the struct APIError 
	switch (resp.StatusCode) {
		case 401:
			return nil, http.StatusUnauthorized, nil
		case 403:
			return nil, http.StatusForbidden, nil
		case 200:
		case 400:
			return nil, http.StatusNotFound, nil
		case 404:
			return nil, http.StatusNotFound, nil
		case 500:
			return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("%d", result["message"]))	
		default:
			return nil, http.StatusInternalServerError, errors.New(fmt.Sprintf("%d", result["message"]))
	}

	return result, http.StatusOK, nil
}