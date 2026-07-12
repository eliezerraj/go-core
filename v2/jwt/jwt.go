package jwt

import (
	"context"
	"crypto/rsa"
	"encoding/base64"
	"fmt"
	"math/big"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog"
)

var childLogger = zerolog.New(os.Stdout).
	With().
	Str("component", "go-core").
	Str("package", "v2.jwt").
	Timestamp().
	Logger()

type JWTValidator struct {
	token     string
	publicKey string
	rsaPublic *rsa.PublicKey
}

func (j *JWTValidator) NewJWTValidator(publicKey string) *JWTValidator {
	childLogger.Debug().
		Str("func", "NewJWTValidator").Send()

	rsaPublic, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		// Fallback to raw JWK modulus (n) with default exponent AQAB.
		rsaPublic, err = parseRSAPublicKeyFromJWK(publicKey, "AQAB")
		if err != nil {
			childLogger.Error().
				Err(err).
				Str("func", "NewJWTValidator").
				Msg("Failed to parse RSA public key")
			return nil
		}
	}

	return &JWTValidator{
		publicKey: publicKey,
		rsaPublic: rsaPublic,
	}
}

func (j *JWTValidator) NewJWTValidatorFromJWK(n, e string) *JWTValidator {
	childLogger.Debug().
		Str("func", "NewJWTValidatorFromJWK").Send()

	rsaPublic, err := parseRSAPublicKeyFromJWK(n, e)
	if err != nil {
		childLogger.Error().
			Err(err).
			Str("func", "NewJWTValidatorFromJWK").
			Msg("Failed to parse JWK RSA public key")
		return nil
	}

	return &JWTValidator{
		publicKey: n,
		rsaPublic: rsaPublic,
	}
}

func parseRSAPublicKeyFromJWK(n, e string) (*rsa.PublicKey, error) {
	nBytes, err := base64.RawURLEncoding.DecodeString(n)
	if err != nil {
		return nil, fmt.Errorf("decode modulus n: %w", err)
	}

	eBytes, err := base64.RawURLEncoding.DecodeString(e)
	if err != nil {
		return nil, fmt.Errorf("decode exponent e: %w", err)
	}

	if len(eBytes) == 0 {
		return nil, fmt.Errorf("invalid exponent e")
	}

	exponent := 0
	for _, b := range eBytes {
		exponent = (exponent << 8) | int(b)
	}

	modulus := new(big.Int).SetBytes(nBytes)
	if modulus.Sign() == 0 || exponent <= 0 {
		return nil, fmt.Errorf("invalid RSA public key values")
	}

	return &rsa.PublicKey{N: modulus, E: exponent}, nil
}

func (j *JWTValidator) ValidateToken(ctx context.Context, tokenString string) (bool, error) {
	childLogger.Debug().
		Str("func", "ValidateToken").Send()

	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return j.rsaPublic, nil
	})

	childLogger.Debug().
			Interface("time.now", time.Now()).
			Interface("claims.ExpiresAt.Time", claims.ExpiresAt.Time).
			Send()

	if err != nil {
		childLogger.Error().
			Err(err).
			Str("func", "ValidateToken").
			Msg("Failed to parse JWT token")
		return false, err
	}

	if claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
		childLogger.Warn().
			Str("func", "ValidateToken").
			Interface("time.now", time.Now()).
			Interface("claims.ExpiresAt.Time", claims.ExpiresAt.Time).
			Msg("JWT token is expired")
		return false, jwt.ErrTokenExpired
	}

	return token.Valid, nil
}
