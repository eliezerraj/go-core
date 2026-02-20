package http

import (
	"os"
	"testing"
	"context"
	"time"
	"github.com/rs/zerolog"
)

var logger = zerolog.New(os.Stdout).
					With().
					Timestamp().
					Str("component","go-core").
					Str("package", "http").
					Logger()

func TestGoCore_HttpService(t *testing.T){

	headers := map[string]string{
		"Content-Type":  "application/json;charset=UTF-8",
		"Authorization": "MY-TOKEN:123",
		"X-Request-Id": "REQ:111-QQQ-RRR-444",
	}

	httpClientParameter := HttpClientParameter {
		Url: 	"http://127.0.0.1:7000/info",
		Method: "GET",
		Timeout: 3,
		Headers: &headers,
	}

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 10 ) * time.Second)
	defer cancel()

	httpService := NewHttpService(&logger)

	//-----------------------------------------------------
	res, statusCode, err := httpService.DoHttp(ctx, httpClientParameter )
	if err != nil {
		t.Errorf("err : %s", err)
	}

	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)
	//--------------------------------

	/*httpClient = HttpClient {
		Url: 	"https://go-XXXXXX-lambda.architecture.stacaradhras.io/oauth_credential",
		Method: "POST",
		Timeout: 3,
		Headers: &headers,
	}
	payload := map[string]string{"user":"admin","password":"admin"}

	res, statusCode, err = apiService.CallRestApi(ctx, httpClient, payload)
	if err != nil {
		t.Errorf("err : %s", err)
		return
	}

	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)

	//--------------------------------

	httpClient = HttpClient {
		Url: 	"https://go-XXXXX-lambda.architecture.caradhras.io/oauth_credential",
		Method: "POST",
		Timeout: 3,
		Headers: &headers,
	}
	payload = map[string]string{"user":"admin--1","password":"admin"}

	res, statusCode, err = apiService.CallRestApi(ctx, httpClient, payload)
	if statusCode == 404 {
		t.Logf("OK the data not existes: %s", err)
		return
	} else {
		t.Errorf("ops...the data supose not exists err : %s", err)
	}
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)*/

	
	//--------------------------------
}
