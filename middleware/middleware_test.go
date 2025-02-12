package middleware

import (
	"testing"
	"net/http"
	//"net/http/httptest"
	"encoding/json"

	"github.com/gorilla/mux"
)

func TestCore_Middleware(t *testing.T){
	var testeToolsCore ToolsMiddleware

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(testeToolsCore.MiddleWareHandlerHeader)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		json.NewEncoder(rw).Encode("appServer")
	})
	
	//httptest.NewRequest("GET", "http://example.com/a", nil)
	//require.NoError(err, "create request")

	//addDebit.HandleFunc("/add", middleware.MiddleWareErrorHandler(httpWorkerAdapter.Add))

	/*timerHandler  := func(w http.ResponseWriter, r *http.Request) {}

	req := httptest.NewRequest(http.MethodGet, "http://example.com", nil)
	res := httptest.NewRecorder()
	timerHandler(res, req)
 
    tim := testeToolsCore.MiddleWareHandlerHeader(timerHandler)
    tim.ServeHTTP(res, req)*/
}