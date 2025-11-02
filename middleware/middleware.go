package middleware

import (
	"context"	
	"net/http"

	"github.com/eliezerraj/go-core/coreJson"
	"github.com/rs/zerolog/log"
)


var childLogger = log.With().Str("component","go-core").Str("package", "middleware").Logger()

var core_json coreJson.CoreJson

type ToolsMiddleware struct {
}

// About middleware http header
func (t *ToolsMiddleware) MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		childLogger.Debug().Msg("................ MiddleWareHandlerHeader. (INICIO) ..........")

		// --- CORS CONFIGURATION ---
		w.Header().Set("Access-Control-Allow-Origin", "*") // or restrict to specific origin
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin")

		// --- Security Headers ---
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "same-origin")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		// --- Set request ID in context if exists ---
		if vals := r.Header.Values("X-Request-Id"); len(vals) > 0 {
			ctx := context.WithValue(r.Context(), "trace-request-id", vals[0])
			r = r.WithContext(ctx)
		}

		childLogger.Debug().Msg("........... MiddleWareHandlerHeader. (FIM) ..........")

		next.ServeHTTP(w, r)
	})
}

type apiFunc func(w http.ResponseWriter, r *http.Request) error

// About middleware http error header
func (t *ToolsMiddleware) MiddleWareErrorHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		childLogger.Debug().Msg("................ MiddleWareErrorHandler (INICIO - RESPONSE/ERROR)  ...............")
		if err := h(w, r); err != nil {
			if e, ok := err.(*coreJson.APIError); ok{
				core_json.WriteJSON(w, e.StatusCode, e)
			}
		}
		childLogger.Debug().Msg(".......... MiddleWareErrorHandler (FIM - RESPONSE/ERROR)  ...............")
	 }
}