package api

import (
	"testing"
	"context"
	"time"
)

func TestGoCore_Api(t *testing.T){

	url := "http://127.0.0.1:5001/info"
	method := "GET"
	header_x_apigw_api_id := "" 
	header_authorization := ""
	header_x_resquest_id := ""

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	apiService := NewRestApiService()

	res, statusCode, err := apiService.CallApi(ctx, url, method, &header_x_apigw_api_id, &header_authorization, &header_x_resquest_id, nil)
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)

	url = "https://go-XXXXXX-lambda.architecture.caradhras.io/oauth_credential"
	method = "POST"
	payload := map[string]string{"user":"admin","password":"admin"}

	res, statusCode, err = apiService.CallApi(ctx, url, method, &header_x_apigw_api_id, &header_authorization,  &header_x_resquest_id,payload)
	if err != nil {
		t.Errorf("err : %s", err)
		return
	}

	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)

	url = "https://go-XXXXX-lambda.architecture.caradhras.io/oauth_credential"
	method = "POST"
	payload = map[string]string{"user":"admin--1","password":"admin"}

	res, statusCode, err = apiService.CallApi(ctx, url, method, &header_x_apigw_api_id, &header_authorization,  &header_x_resquest_id,payload)
	if statusCode == 404 {
		t.Logf("OK the data not existes: %s", err)
		return
	} else {
		t.Errorf("ops...the data supose not exists err : %s", err)
	}
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)
}