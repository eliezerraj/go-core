package coreJson

import (
	"testing"

	"bytes"
	"net/http"
	"net/http/httptest"
)

var testJson = []struct {
	payload          string
}{
	{payload: `{"foo": "bar","asb":"123"}`},
}

func TestCoreJson_ReadJSON(t *testing.T){
	var coreJson CoreJson
	
	req, err := http.NewRequest("POST", "/", bytes.NewReader([]byte(testJson[0].payload)))
	if err != nil {
		t.Log("Error:", err)
	}
	r := httptest.NewRecorder()

	var decodedJSON struct {
		Foo string `json:"foo"`
	}

	err = coreJson.ReadJSON(req, r, &decodedJSON)
	if err != nil {
		t.Errorf("err : %s:", err)
	}

	req.Body.Close()
}

func TestCoreJson_WriteJSON(t *testing.T){
	var coreJson CoreJson

	w := httptest.NewRecorder()
	payload := JSONResponse{
		Error:   false,
		Message: "test WriteJSON",
	}
	
	err := coreJson.WriteJSON(w, http.StatusOK, payload)
	if err != nil {
		t.Errorf("failed to write JSON: %v", err)
	}

	t.Logf("w: %v", w)
}