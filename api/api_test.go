package api

import (
	"testing"
	"context"
	"time"
)

/*func TestGoCore_CallRestApi(t *testing.T){

	headers := map[string]string{
		"Content-Type":  "application/json;charset=UTF-8",
		"Authorization": "MY-TOKEN:123",
		"X-Request-Id": "REQ:111-QQQ-RRR-444",
	}

	httpClient := HttpClient {
		Url: 	"http://127.0.0.1:5001/info",
		Method: "GET",
		Timeout: 3,
		Headers: &headers,
	}

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	apiService := NewRestApiService()

	//-----------------------------------------------------
	res, statusCode, err := apiService.CallRestApi(ctx, httpClient, nil)
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)
	//--------------------------------

	httpClient = HttpClient {
		Url: 	"https://go-XXXXXX-lambda.architecture.caradhras.io/oauth_credential",
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
	t.Logf("res: %v", res)

	
	//--------------------------------
}*/

func TestGoCore_CallRestApi2(t *testing.T){

	headers := map[string]string{
		"Content-Type":  "application/json;charset=UTF-8",
		"Host": "go-limit.XXXXX.XXXXXX.io",
	}

	httpClient := HttpClient {
		Url: 	"https://nlb-eks-arch-02.XXXXXX.XXXXXX.io/info",
		Method: "GET",
		Timeout: 15,
		Headers: &headers,
	}

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	apiService := NewRestApiService()

	//-----------------------------------------------------
	res, statusCode, err := apiService.CallRestApi(ctx, httpClient, nil)
	if err != nil {
		t.Errorf("err : %s", err)
	}
	if statusCode == 200 {
		t.Logf("OK the data existes: %s", err)
		return
	} else {
		t.Errorf("ops...the data supose exists err : %s", err)
	}
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)	
}