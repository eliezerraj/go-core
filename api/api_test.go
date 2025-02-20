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

	ctx, cancel := context.WithTimeout(	context.Background(), 
										time.Duration( 30 ) * time.Second)
	defer cancel()

	apiService := NewRestApiService()

	res, statusCode, err := apiService.CallApi(ctx, url, method, &header_x_apigw_api_id, &header_authorization, nil)
	if err != nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)

	url = "https://go-login.architecture.caradhras.io/loginRSA"
	method = "POST"
	payload := map[string]string{"user":"admin","password":"admin"}

	res, statusCode, err = apiService.CallApi(ctx, url, method, &header_x_apigw_api_id, &header_authorization, payload)
	if err != nil {
		t.Errorf("err : %s", err)
	}

	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)

	url = "https://go-login.architecture.caradhras.io/loginRSA"
	method = "POST"
	payload = map[string]string{"user":"admin1","password":"admin"}

	res, statusCode, err = apiService.CallApi(ctx, url, method, &header_x_apigw_api_id, &header_authorization, payload)
	if err == nil {
		t.Errorf("err : %s", err)
	}
	t.Logf("err : %s", err)
	t.Logf("statusCode: %v", statusCode)
	t.Logf("res: %v", res)
}