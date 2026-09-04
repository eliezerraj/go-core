package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"crypto"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/lestrrat-go/jwx/v2/jwk"
 	"github.com/golang-jwt/jwt/v5"

	"github.com/eliezerraj/go-core/v3/logger"
	"github.com/eliezerraj/go-core/v3/httpclient"
)

type AuthService struct {
	httpConfig *httpclient.HttpConfig
	jwksUrl string
	dryRun bool
	header string
	publicKeys map[string]crypto.PublicKey
	httpTimeout time.Duration
}

const (
	AcceptHeader      = "Accept"
	ContentTypeHeader = "Content-Type"
	ConnectionHeader  = "Connection"
	KeepAlive         = "keep-alive"
)

func NewAuthService(jwksUrl string, dryRun bool, header string, httpTimeout time.Duration) *AuthService {
	logger.Info(context.Background(), "Initializing AuthService SUCCESSFULLY")

	httpConfig := &httpclient.HttpConfig{
		Timeout:             5 * time.Second,
		KeepAlive:           5 * time.Second,
		IdleConnTimeout:     5 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
		MaxConnsPerHost:     100,
		ServiceName:         "auth-service",
	}

	return &AuthService{
		httpConfig:  httpConfig,
		jwksUrl:     jwksUrl,
		dryRun:      dryRun,
		header:      header,
		httpTimeout: httpTimeout,
	}
}

// GetJwksUrl retrieves the JWKS URL from the auth service and ensures it is accessible.
func(a *AuthService) GetJwksUrl(ctx context.Context) ( error) {
	logger.Info(ctx, "Retrieving JWKS URL from auth service")

	ctxHttpTimeout, cancel := context.WithTimeout(ctx, a.httpTimeout)
	defer cancel()
	
	authHttpClient := httpclient.NewHttpClient(a.httpConfig)

	method := "GET"
	req, err := http.NewRequestWithContext(ctxHttpTimeout, method, a.jwksUrl, nil)
	if err != nil {
		logger.Error(ctx, "Failed to create request", zap.Error(err))
		return err
	}

	headers := map[string]string{
		ConnectionHeader:  KeepAlive,
		AcceptHeader:      "application/json",
		ContentTypeHeader: "application/json",
		KeepAlive: "timeout=5, max=1000",
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := authHttpClient.Do(req.WithContext(ctx))
	if err != nil {
		logger.Error(ctx, "Failed to perform request", zap.Error(err))
		return err
	}
	defer resp.Body.Close()
		
	switch resp.StatusCode {
	case http.StatusOK:
		// Continue processing
	case http.StatusNotFound:
		logger.Error(ctx, "Auth service returned 404 Not Found")
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	default:
		logger.Error(ctx, "Auth service returned unexpected status", zap.Int("status", resp.StatusCode))
		return fmt.Errorf("auth service returned status: %d", resp.StatusCode)
	}

	// Decode the response body into a WellKnownJwks struct
	var res_wellknowjws WellKnownJwks
	if err := json.NewDecoder(resp.Body).Decode(&res_wellknowjws); err != nil {
		logger.Error(ctx, "Failed to decode response", zap.Error(err))
		return err
	}

   // 1. Convert the JWK struct into []byte JSON
    keyBytes, err := json.Marshal(res_wellknowjws.Keys[0])
    if err != nil {
        logger.Error(ctx, "Failed to marshal JWK struct", zap.Error(err))
        return err
    }

    // 2. Parse []byte
    publicParseKey, err := jwk.ParseKey(keyBytes)
    if err != nil {
        logger.Error(ctx, "Failed to parse JWK", zap.Error(err))
        return err
    }

    var rawKey crypto.PublicKey
    if err := publicParseKey.Raw(&rawKey); err != nil {
        logger.Error(ctx, "Failed to get raw key from JWK", zap.Error(err))
        return err
    }

    // 3. Store in map: map[string]crypto.PublicKey
    if a.publicKeys == nil {
        a.publicKeys = make(map[string]crypto.PublicKey)
    }

    kid := res_wellknowjws.Keys[0].KeyID
    a.publicKeys[kid] = rawKey

	return nil
}

// VerifyToken verifies the given JWT token using the public keys stored in the AuthService. It returns an error if the token is invalid or cannot be verified.
func(a *AuthService) VerifyToken(ctx context.Context, token string) (*Claims, error) {
	logger.Info(ctx, "Verifying token")

    parsedToken, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			logger.Error(ctx, "Unexpected signing method", zap.Any("alg", token.Header["alg"]))
            return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
        }

		// Extract the 'kid' from the token header and use it to look up the corresponding public key.
        kid, ok := token.Header["kid"].(string)
        if !ok {
            return nil, fmt.Errorf("missing or invalid 'kid' header in token")
        }

		// Look up the public key KID corresponding to the 'kid' in the token header.
        pubKey, ok := a.publicKeys[kid]
        if !ok {
            return nil, fmt.Errorf("public key not found for kid: %s", kid)
        }

        return pubKey, nil
    })
	
    if err != nil {
        logger.Error(ctx, "Failed to parse token", zap.Error(err))
        return nil, err
    }

    if !parsedToken.Valid {
        logger.Error(ctx, "Token is invalid")
        return nil, fmt.Errorf("token is invalid")
    }

	claims, ok := parsedToken.Claims.(*Claims)
	if !ok {
		logger.Error(ctx, "Failed to cast token claims")
		return nil, fmt.Errorf("failed to cast token claims")
	}
	
    token = parsedToken.Raw

	return claims, nil
}