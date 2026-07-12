package dpop

import (
	"log"
	"fmt"
	"time"
	"math"
	"errors"
	"encoding/base64"
	"encoding/json"
	"crypto/ecdsa"
	"crypto/elliptic"
	"math/big"
	"crypto/sha256"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("your-very-secure-secret-key")

type Confirmation struct {
	Jkt string `json:"jkt"`
}

type CustomClaims struct {
	Username string       `json:"username"`
	Scopes   []string     `json:"scopes"`
	Cnf      *Confirmation `json:"cnf,omitempty"`
	jwt.RegisteredClaims
}

type BearerToken struct {
	Token     	string    `json:"token"`
}

func AuthorizationServer(username string,
						 dpopToken string,
						 method string,
						 url string) (*BearerToken, error) {

	log.Printf("\n AuthorizationServer: username=%s, dpopToken=%s, method=%s, url=%s", username, dpopToken, method, url)

	var jkt string

	// Extract jwt from Dpop
	parsed, err := jwt.Parse(dpopToken, func(t *jwt.Token) (interface{}, error) {
		jwkMap, ok := t.Header["jwk"].(map[string]any)
		if !ok {
			return nil, errors.New("jwk header missing or invalid")
		}

		kty, _ := jwkMap["kty"].(string)
		crv, _ := jwkMap["crv"].(string)
		xStr, _ := jwkMap["x"].(string)
		yStr, _ := jwkMap["y"].(string)

		if kty != "EC" || crv != "P-256" {
			return nil, errors.New("unsupported JWK key type or curve")
		}

		xBytes, err := base64.RawURLEncoding.DecodeString(xStr)
		if err != nil {
			return nil, err
		}
		yBytes, err := base64.RawURLEncoding.DecodeString(yStr)
		if err != nil {
			return nil, err
		}

		publicKey := &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}

		// Compute JWK Thumbprint (RFC 7638)
		// The fields MUST be in lexicographical order: crv, kty, x, y
		thumbprintMap := map[string]string{
			"crv": crv,
			"kty": kty,
			"x":   xStr,
			"y":   yStr,
		}
		
		jwkJson, err := json.Marshal(thumbprintMap)
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256(jwkJson)
		jkt = base64.RawURLEncoding.EncodeToString(hash[:])
		
		log.Printf("JWK Thumbprint (jkt): %s", jkt)

		return publicKey, nil
	})
    if err != nil {
        return nil, err
    }

	// Validate the DPoP token claims
	claimsDpop := parsed.Claims.(jwt.MapClaims)
    if claimsDpop["htm"] != method {
        return nil, errors.New("invalid method")
    }

    if claimsDpop["htu"] != url {
        return nil, errors.New("invalid url")
    }

	iat := int64(claimsDpop["iat"].(float64))

    if math.Abs(float64(time.Now().Unix()-iat)) > 60 {
        return nil,errors.New("token expired")
    }

	var scopes = []string{"read", "write", "delete"}
	// 3. Populate the claims (Registered + Custom)
	claims := CustomClaims{
		Username: username,
		Scopes:   scopes,
		Cnf: &Confirmation{
			Jkt: jkt,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)), // Expires in 24 hours
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "testgocore.auth",
			Subject:   "user credentials",
		},
	}

	// 4. Create the token using the claims and a signing method
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 5. Sign the token with your secret key
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		fmt.Printf("Error signing token: %v\n", err)
		return nil, err
	}

	// 6. Format as a standard Authorization Header Bearer token
	// NOTE: Return the raw token string so it can be parsed correctly, or use tokenString in validate
	bearerHeader := fmt.Sprintf("Bearer %s", tokenString)

	// Compute the SHA-256 hash of the DPoP token

	return &BearerToken{
		Token: bearerHeader,
	}, nil
}

func AuthorizationServerTokenDPopValidation(dpopToken string, bearerToken string, method string, url string) error {
	log.Printf("\n AuthorizationServerTokenDPopValidation:\n dpopToken=%s\n bearerToken=%s\n method=%s, url=%s", dpopToken, bearerToken, method, url)

	// Validate Bearer Token
	// The bearerToken might be formatted as "Bearer <token>", so we need to strip the prefix
	cleanToken := bearerToken
	if len(bearerToken) > 7 && bearerToken[:7] == "Bearer " {
		cleanToken = bearerToken[7:]
	}

	parsedBearer, err := jwt.ParseWithClaims(cleanToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})	
	if err != nil {
		return err
	}

	//Validate Bearer Token Claims
	claimsBearer := parsedBearer.Claims.(*CustomClaims)
	ExpiresAt := claimsBearer.RegisteredClaims.ExpiresAt
	if ExpiresAt == nil || time.Now().After(ExpiresAt.Time) {
		return errors.New("bearer token expired")
	}

	// Extract jkt from Bearer Token
	jkt := claimsBearer.Cnf.Jkt
	log.Printf("Extracted jkt from Bearer Token: %s", jkt)

	// Validate DPoP Token
	var jkt_DPoP string
	parsed, err := jwt.Parse(dpopToken, func(t *jwt.Token) (interface{}, error) {
		jwkMap, ok := t.Header["jwk"].(map[string]any)
		if !ok {
			return nil, errors.New("jwk header missing or invalid")
		}

		kty, _ := jwkMap["kty"].(string)
		crv, _ := jwkMap["crv"].(string)
		xStr, _ := jwkMap["x"].(string)
		yStr, _ := jwkMap["y"].(string)

		if kty != "EC" || crv != "P-256" {
			return nil, errors.New("unsupported JWK key type or curve")
		}

		xBytes, err := base64.RawURLEncoding.DecodeString(xStr)
		if err != nil {
			return nil, err
		}
		yBytes, err := base64.RawURLEncoding.DecodeString(yStr)
		if err != nil {
			return nil, err
		}

		publicKey := &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xBytes),
			Y:     new(big.Int).SetBytes(yBytes),
		}

		// Compute JWK Thumbprint (RFC 7638)
		// The fields MUST be in lexicographical order: crv, kty, x, y
		thumbprintMap := map[string]string{
			"crv": crv,
			"kty": kty,
			"x":   xStr,
			"y":   yStr,
		}
		
		jwkJson, err := json.Marshal(thumbprintMap)
		if err != nil {
			return nil, err
		}

		hash := sha256.Sum256(jwkJson)
		jkt_DPoP = base64.RawURLEncoding.EncodeToString(hash[:])
		
		log.Printf("JWK Thumbprint (jkt_DPoP): %s", jkt_DPoP)

		return publicKey, nil
	})
    if err != nil {
        return err
    }

	// Validate the DPoP token claims
	claimsDpop := parsed.Claims.(jwt.MapClaims)
    if claimsDpop["htm"] != method {
        return errors.New("invalid method")
    }

    if claimsDpop["htu"] != url {
        return errors.New("invalid url")
    }

	iat := int64(claimsDpop["iat"].(float64))

    if math.Abs(float64(time.Now().Unix()-iat)) > 60 {
        return errors.New("token expired")
    }

	// Validate the DPoP signature and claims
	log.Printf("---------------------------------")
	log.Printf("JWK Thumbprint (jkt_Bearer): %s", jkt)
	log.Printf("JWK Thumbprint (jkt_DPoP)  : %s", jkt_DPoP)
	log.Printf("---------------------------------")

	// Implement the DPoP token validation logic here
	return nil
}