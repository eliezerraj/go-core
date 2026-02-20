package middleware

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"

	"github.com/rs/zerolog"
	"github.com/google/uuid"
)

// contextKey is a type for context keys to avoid collisions
type contextKey string
const RequestIDKey contextKey = "request-id"

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

type MiddleWare struct {
	logger    zerolog.Logger
	corsConfig *CORSConfig
}

// APIError represents API error response
type APIError struct {
	StatusCode int    `json:"status_code"`
	Msg        string `json:"message"`
	RequestID  string `json:"request_id,omitempty"`
}

// NewMiddleWare creates a new middleware instance
func NewMiddleWare(appLogger *zerolog.Logger) *MiddleWare {
	logger := appLogger.With().
						Str("component", "go-core.v2.middleware").
						Logger()

	return &MiddleWare{
		logger: logger,
		corsConfig: &CORSConfig{
			AllowedOrigins: []string{"*"},
			AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders: []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Origin", "X-Amz-Date", "X-Api-Key", "X-Amz-Security-Token"},
		},
	}
}

// NewMiddleWareWithCORS creates middleware with custom CORS configuration
func NewMiddleWareWithCORS(appLogger *zerolog.Logger, cors *CORSConfig) *MiddleWare {
	logger := appLogger.With().
						Str("component", "go-core.v2.middleware").
						Logger()

	return &MiddleWare{
		logger: logger,
		corsConfig: cors,
	}
}

// GetRequestID retrieves the request ID from context
// Returns empty string if not found
func GetRequestID(ctx context.Context) string {
	id, ok := ctx.Value(RequestIDKey).(string)
	if !ok {
		// Request ID not found in context
		return ""
	}
	return id
}

// DebugContextValues logs context value keys for debugging
// This is a helper function to troubleshoot context propagation issues
func DebugContextValues(ctx context.Context, logger *zerolog.Logger) {
	requestID := GetRequestID(ctx)
	logger.Debug().
		Str("request_id_from_key", requestID).
		Interface("context_keys_available", "Use GetRequestID() function").
		Msg("Context debug info")
}

// WriteJSON writes data as JSON response with given status code
func (m *MiddleWare) WriteJSON(w http.ResponseWriter, code int, data interface{}) error {
	data_json, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal JSON response: %w", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	_, err = w.Write(data_json)
	if err != nil {
		m.logger.Error().
			Err(err).
			Msg("Failed to write HTTP response")
		return fmt.Errorf("failed to write response: %w", err)
	}
	return nil
}

// MiddleWareHandlerHeader sets security and CORS headers, and generates request ID
func (m *MiddleWare) MiddleWareHandlerHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// --- CORS Headers ---
		cors := m.corsConfig
		w.Header().Set("Access-Control-Allow-Origin", cors.AllowedOrigins[0])
		w.Header().Set("Access-Control-Allow-Methods", joinStrings(cors.AllowedMethods))
		w.Header().Set("Access-Control-Allow-Headers", joinStrings(cors.AllowedHeaders))

		// --- Security Headers ---
		w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
		w.Header().Set("Content-Security-Policy", "default-src 'none'; img-src 'self'; script-src 'self'; style-src 'self'; object-src 'none'; frame-ancestors 'none'")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "same-origin")

		// --- Common Headers ---
		w.Header().Set("Content-Type", "application/json")

		// --- Handle CORS preflight ---
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// --- Set request ID in context ---
		requestID := getOrGenerateRequestID(r)
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)
		r = r.WithContext(ctx)

		m.logger.Debug().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Msg("Request ID set in context")

		next.ServeHTTP(w, r)
	})
}

// getOrGenerateRequestID retrieves request ID from header or generates new one
func getOrGenerateRequestID(r *http.Request) string {
	if vals := r.Header.Values("X-Request-Id"); len(vals) > 0 {
		return vals[0]
	}
	return uuid.New().String()
}

// joinStrings joins string slice with comma separator
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += ", "
		}
		result += s
	}
	return result
}

func (e *APIError) Error() string {
	return e.Msg
}

// NewAPIError creates a new API error response
func NewAPIError(err error, requestID string, status ...int) *APIError {
	// Set default status code
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	return &APIError{
		StatusCode: statusCode,
		Msg:        err.Error(),
		RequestID:  requestID,
	}
}

// apiFunc is a handler function that can return an error
type apiFunc func(w http.ResponseWriter, r *http.Request) error

// MiddleWareErrorHandler wraps handler functions and handles their errors
func (m *MiddleWare) MiddleWareErrorHandler(h apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h(w, r); err != nil {
			requestID := GetRequestID(r.Context())

			// Handle APIError types
			if apiErr, ok := err.(*APIError); ok {
				m.logger.Warn().
					Err(err).
					Str("request_id", requestID).
					Int("status_code", apiErr.StatusCode).
					Msg("API error")
				m.WriteJSON(w, apiErr.StatusCode, apiErr)
				return
			}

			// Handle all other error types
			m.logger.Error().
				Err(err).
				Str("request_id", requestID).
				Msg("Internal server error")

			apiErr := NewAPIError(
				fmt.Errorf("internal server error"),
				requestID,
				http.StatusInternalServerError)
			m.WriteJSON(w, apiErr.StatusCode, apiErr)
		}
	}
}

// MiddleWareRecovery recovers from handler panics and returns error response
func (m *MiddleWare) MiddleWareRecovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				requestID := GetRequestID(r.Context())
				m.logger.Error().
					Interface("panic", err).
					Str("request_id", requestID).
					Msg("Handler panic recovered")

				apiErr := NewAPIError(
					fmt.Errorf("internal server error"),
					requestID,
					http.StatusInternalServerError)
				m.WriteJSON(w, apiErr.StatusCode, apiErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}