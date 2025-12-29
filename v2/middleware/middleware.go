package middleware

import (
	"context"	
	"net/http"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/google/uuid"
)

type MiddleWare struct {
	logger 		*zerolog.Logger
}

// About create a new middlware
func NewMiddleWare(appLogger *zerolog.Logger) *MiddleWare {
	logger := appLogger.With().
						Str("component", "go-core.v2.middleware").
						Logger()
	logger.Debug().
			Str("func","NewMiddleWare").Send()

	return &MiddleWare {
		logger: &logger,
	}
}

// About convert JSON
func (t *MiddleWare) WriteJSON(	w http.ResponseWriter, 
								code int, 
								data interface{}) error {
	
	data_json, err := json.Marshal(data)
	if err != nil {
		return err
	}
	
	w.WriteHeader(code)

	_, err = w.Write(data_json)
	if err != nil {
		return err
	}
	return nil
}

// ------------------ MiddleWare Headet ---------------------
func (t *MiddleWare) MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//t.logger.Debug().
		//    	 Msg("................ MiddleWareHandlerHeader. (INICIO) ..........")

		// --- CORS ---
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, Accept, Origin, X-Amz-Date, X-Api-Key, X-Amz-Security-Token")

		// --- Security Headers ---
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "same-origin")

		// --- Common Headers ---
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// --- Set request ID in context if exists ---
		if vals := r.Header.Values("X-Request-Id"); len(vals) > 0 {
			ctx := context.WithValue(r.Context(), "request-id", vals[0])
			r = r.WithContext(ctx)
		} else {
			//t.logger.Debug().
			//		Msg("creating new request-id from uuid")
			ctx := context.WithValue(r.Context(), "request-id", uuid.New().String())
			r = r.WithContext(ctx)
		}
		
		//t.logger.Debug().
		// 		 Msg("........... MiddleWareHandlerHeader. (FIM) ..........")

		next.ServeHTTP(w, r)
	})
}

// ------------------- ERROR --------------------------
type APIError struct {
	StatusCode	int    `json:"status_code"`
	Msg			string `json:"message"`
	TraceId		string `json:"request_id,omitempty"`
}

func (e *APIError) Error() string {
	return e.Msg
}

func (e *APIError) NewAPIError(	err error, 
								traceId string, 
								status ...int) APIError {

	// set a default status code
	statusCode := http.StatusBadRequest
	
	if len(status) > 0 {
		statusCode = status[0]
	}

	return APIError{
		StatusCode: statusCode,
		Msg:		err.Error(),
		TraceId:	traceId,
	}
}

// ------------------ MiddleWare Response ---------------------------
type apiFunc func(w http.ResponseWriter, r *http.Request) error

// About middleware http error header
func (t *MiddleWare) MiddleWareErrorHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		t.logger.Debug().
				 Msg("................ MiddleWareErrorHandler v2 (INICIO - RESPONSE/ERROR)  ...............")
		
		if err := h(w, r); err != nil {
			if e, ok := err.(*APIError); ok{
				t.WriteJSON(w, e.StatusCode, e)
			}
		}
		w.Header().Set("Content-Type", "application/json")
		
		t.logger.Debug().
				 Msg(".......... MiddleWareErrorHandler v2 (FIM - RESPONSE/ERROR)  ...............")
	 }
}