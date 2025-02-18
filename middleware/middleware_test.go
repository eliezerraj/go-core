package middleware

import (
	"testing"
	"net/http"
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
}