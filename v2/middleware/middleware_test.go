package middleware

import (
	"os"
	"testing"
	"net/http"
	"encoding/json"
	"github.com/rs/zerolog"

	"github.com/gorilla/mux"
)

func TestCore_Middleware(t *testing.T){
	
	var logger = zerolog.New(os.Stdout).
						With().
						Str("component", "TestCore_Middleware").
						Logger()

	testeToolsCore := NewMiddleWare(&logger)

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.Use(testeToolsCore.MiddleWareHandlerHeader)

	myRouter.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		json.NewEncoder(rw).Encode("appServer")
	})
}