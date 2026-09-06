package auth

import (
	"context"
	"fmt"
	"net/http"
	"encoding/json"
	"bytes"
	"time"
	"sync"

	"go.uber.org/zap"
	"github.com/eliezerraj/go-core/v3/logger"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type AuthClientServiceOption func(*AuthClientService)

type AuthClientService struct {
	DryRun			bool
	RefreshInterval	int
	AuthURL			string
	RefreshURL		string
	ClientID    	string 
	ClientSecret	string 
	HttpClient		*http.Client

	accessToken		string `json:"access_token"`

	mu     sync.RWMutex
	workerTokenOnce sync.Once
	expiresAt   time.Time

	ctx            context.Context
	cancel         context.CancelFunc
}

type authRequest struct {
    ClientID     string `json:"client_id"`
    ClientSecret string `json:"client_secret"`
}

type authorizationResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

// NewAuthClientService initializes a new instance of AuthClientService with the provided options.
func NewAuthClientService(options ...AuthClientServiceOption) *AuthClientService {
	logger.InfoOutCtx("Initializing AuthClientService SUCCESSFULLY")

	ctx, cancel := context.WithCancel(context.Background())

	// Initialize the HTTP transport for the client with OpenTelemetry instrumentation
    transport := &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        MaxConnsPerHost:     100,
        IdleConnTimeout:     90 * time.Second,
    }

	// Create the HTTP client with the instrumented transport
	httpClient := &http.Client{
		Transport: otelhttp.NewTransport(transport),
		Timeout:   15 * time.Second,
	}
	
	// Initialize the AuthClientService with default values and the HTTP client
	authClientService := &AuthClientService{
		DryRun:         false,
		HttpClient:     httpClient,
		ctx:            ctx,
		cancel:         cancel,
	}

	// Apply the provided options to the AuthClientService instance
	for _, option := range options {
		option(authClientService)
	}

	return authClientService
}

// GetToken returns the current cached bearer token
func (a *AuthClientService) GetToken() string {
    a.mu.RLock()
    defer a.mu.RUnlock()
    return a.accessToken
}

// Close terminates the AuthClientService by canceling its context.
func (a *AuthClientService) Close() {
    a.cancel()
}

// StartAuthenticate initiates the authentication process and returns the access token. If dry run is enabled, it skips the authentication and returns an empty string.
func (a *AuthClientService) StartAuthenticate(ctx context.Context) (string, error) {
	logger.Info(ctx, "Starting authentication process")

	tracer := otel.Tracer("authClientService")
	ctx, span := tracer.Start(ctx, "go-core.auth.AuthClientService.Authenticate", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	if a.DryRun == true {
		logger.Info(ctx, "Dry run enabled, skipping authentication")
		return "", nil
	}

	var startErr error
	a.workerTokenOnce.Do(func() {
		accessToken, err := a.DoHttpRequest(ctx)
		if err != nil {
			logger.Error(ctx, "Failed to fetch initial token", zap.Error(err))
			return
		}

		a.mu.Lock()
		a.accessToken = accessToken.AccessToken

		a.expiresAt = time.Now().Add(time.Duration(accessToken.ExpiresIn) * time.Second)
		a.mu.Unlock() 

		logger.Info(ctx, "a.accessToken / a.expiresAt", zap.String("accessToken", a.accessToken), zap.Time("expiresAt", a.expiresAt))
		
		go a.TokenWorker()
	})

	if startErr != nil {
		return "", startErr
	}

	return a.GetToken(), nil
}

// TokenWorker is responsible for managing the token refresh process. It runs in a separate goroutine and refreshes the token when necessary.
func (a *AuthClientService) TokenWorker() {
	logger.Info(a.ctx, "Starting token refresh process")

	for {
        a.mu.RLock()
        expiresAt := a.expiresAt
        a.mu.RUnlock()

        // Refresh when 80% of TTL is reached (or 30s before expiration)
        refreshTime := expiresAt.Add(-30 * time.Second)
        sleepDuration := time.Until(refreshTime)

        if sleepDuration <= 0 {
            sleepDuration = 5 * time.Second
        }

        select {
        case <-a.ctx.Done():
            return // graceful shutdown
        case <-time.After(sleepDuration):
        }

        accessToken, err := a.DoHttpRequest(a.ctx)
        if err != nil {
            logger.Error(a.ctx, "Failed to refresh token, retrying in 5s", zap.Error(err))
            time.Sleep(5 * time.Second)
            continue
        }

        a.mu.Lock()
        a.accessToken = accessToken.AccessToken
        a.expiresAt = time.Now().Add(time.Duration(accessToken.ExpiresIn) * time.Second)
        a.mu.Unlock()

        logger.Info(a.ctx, "Token refreshed successfully", zap.String("accessToken", a.accessToken), zap.Time("expiresAt", a.expiresAt))
	}
}

// DoHttpRequest performs the HTTP request to the authentication service to obtain a new access token.
func (a *AuthClientService) DoHttpRequest(ctx context.Context) (*authorizationResponse, error) {
	logger.Info(ctx, "Do http request to auth service")

	payload := authRequest{
        ClientID:     a.ClientID,
        ClientSecret: a.ClientSecret,
    }

	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error(ctx, "Failed to marshal auth client request", zap.Error(err))
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, a.AuthURL, bytes.NewReader(body))
	if err != nil {
		logger.Error(ctx, "Failed to create request", zap.Error(err))
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := a.HttpClient.Do(req)
	if err != nil {
		logger.Error(ctx, "Failed to execute request", zap.Error(err))
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		// Continue processing
	case http.StatusCreated:
		// Continue processing
	case http.StatusNotFound:
		logger.Error(ctx, "Auth service returned 404 Not Found")
		return nil, fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	default:
		logger.Error(ctx, "Auth service returned unexpected status", zap.Int("status", resp.StatusCode))
		return nil, fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	var result authorizationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Error(ctx, "Failed to decode response", zap.Error(err))
		return nil, err
	}

	// If the expires_in value is not come from the auth service, set a default value
	if result.ExpiresIn <= 0 {
        result.ExpiresIn = 300// default 5 minutes
    }

	if a.RefreshInterval > 0 {
		result.ExpiresIn = a.RefreshInterval // set to the configured refresh interval
	}

	return &result, nil
}