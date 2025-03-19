package middleware

import (
	"fmt"
	"context"	
	"net/http"

	"github.com/eliezerraj/go-core/coreJson"
	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("go-core", "middleware").Logger()
var core_json coreJson.CoreJson

type ToolsMiddleware struct {
}

// About middleware http header
func (t *ToolsMiddleware) MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childLogger.Debug().Msg("-------------- MiddleWareHandlerHeader. (INICIO)  --------------")

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,access-control-allow-origin, access-control-allow-headers")
		w.Header().Set("strict-transport-security","max-age=63072000; includeSubdomains; preloa")
		w.Header().Set("content-security-policy","default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")
		w.Header().Set("x-content-type-option","nosniff")
		w.Header().Set("x-frame-options","DENY")
		w.Header().Set("x-xss-protection","1; mode=block")
		w.Header().Set("referrer-policy","same-origin")
		w.Header().Set("permission-policy","Content-Type,access-control-allow-origin, access-control-allow-headers")

		// set the requet-id
		ctx := context.WithValue(r.Context(), "trace-request-id", fmt.Sprintf("%v",r.Header["X-Request-Id"]))
		r = r.WithContext(ctx)
	
		childLogger.Debug().Msg("-------------- MiddleWareHandlerHeader. (FIM) ----------------")

		next.ServeHTTP(w, r)
	})
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// About middleware http error header
func (t *ToolsMiddleware) MiddleWareErrorHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childLogger.Debug().Msg("------- MiddleWareErrorHandler (INICIO - ERROR)  --------------")
		if err := h(w, r); err != nil {
			if e, ok := err.(*coreJson.APIError); ok{
				core_json.WriteJSON(w, e.StatusCode, e)
			}
		}
		childLogger.Debug().Msg("------- MiddleWareErrorHandler (FIM - ERROR)  --------------")
	 }
}