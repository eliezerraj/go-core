package dpop

import (
	"time"
	"math"
	"math/big"
	"github.com/golang-jwt/jwt/v5"
	"errors"
	"crypto/rand"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"github.com/rs/zerolog"
	"github.com/google/uuid"
)

type DPoPClaims struct {
    Htm string `json:"htm"`
    Htu string `json:"htu"`
    Iat int64  `json:"iat"`
    Jti string `json:"jti"`
    Ath string `json:"ath,omitempty"`
}

type KeyPEM struct {
	privPEM string `json:"private_pem"`
	pubPEM  string `json:"public_pem"`
}

type DPoP struct {
	token string
	logger *zerolog.Logger
}

// CreateKeys generates a new ECDSA P-256 key pair and returns the private and public keys in PEM format.
func (d *DPoP) CreateKeys() (*KeyPEM, error) {
	d.logger.Debug().
		Str("func","CreateKeys").Send()

	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return nil, err
	}

	// Marshal ECDSA Private Key
	privBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return nil, err
	}
	privPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: privBytes})

	// Marshal ECDSA Public Key
	pubBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return nil, err
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes})

	return &KeyPEM{
		privPEM: string(privPEM),
		pubPEM:  string(pubPEM),
	}, nil
}

// NewDPoP creates a new instance of DPoP with the provided logger.
func NewDPoP(logger *zerolog.Logger,) *DPoP {
	logger.Debug().
		Str("func","NewDPoP").Send()
			
	return &DPoP{
		logger: logger,
	}
}

// Create a jwt Dpop token without access token
func (d *DPoP) CreateTokenDpopNoAuth(htm, htu string, privPEM string, pubPEM string) (string, error) {
	d.logger.Debug().
		Str("func","CreateTokenDpopNoAuth").Send()
		
	claims := DPoPClaims{
		Htm: htm,
		Htu: htu,
		Iat: time.Now().Unix(),
		Jti: uuid.NewString(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"htm": claims.Htm,
		"htu": claims.Htu,
		"iat": claims.Iat,
		"jti": claims.Jti,
	})

	block, _ := pem.Decode([]byte(privPEM))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		d.logger.Error().
			Str("func", "CreateTokenDpopNoAuth").
			Msg("failed to decode PEM block containing EC private key")
		return "", errors.New("failed to decode PEM block containing EC private key")
	}

	blockPub, _ := pem.Decode([]byte(pubPEM))
	if blockPub == nil || blockPub.Type != "PUBLIC KEY" {
		d.logger.Error().
			Str("func", "CreateTokenDpopNoAuth").
			Msg("failed to decode PEM block containing EC public key")
		return "", errors.New("failed to decode PEM block containing EC public key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return "", err
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
	if err != nil {
		d.logger.Error().Err(err).Msg("failed to parse public key")
		return "", err
	}
	publicKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("not an ECDSA public key")
	}

	token.Header["typ"] = "dpop+jwt"

	// JWK must contain raw coordinate bytes encoded as Base64URL
	x := publicKey.X.Bytes()
	y := publicKey.Y.Bytes()

	// Ensure padding to the curve's coordinate length (32 bytes for P-256)
	byteLen := (publicKey.Curve.Params().BitSize + 7) / 8
	if len(x) < byteLen {
		padding := make([]byte, byteLen-len(x))
		x = append(padding, x...)
	}
	if len(y) < byteLen {
		padding := make([]byte, byteLen-len(y))
		y = append(padding, y...)
	}

	jwk := map[string]any{
		"kty": "EC",
		"crv": "P-256",
		"x":   base64.RawURLEncoding.EncodeToString(x),
		"y":   base64.RawURLEncoding.EncodeToString(y),
	}

	token.Header["jwk"] = jwk

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return "", err
	}

	return signedToken, nil
}

func (d *DPoP) CreateTokenDPopWithAccessToken(	htm string, 
												htu string, 
												accessToken string, 
												privPEM string, 
												pubPEM string) (string, error) {

	hash := sha256.Sum256([]byte(accessToken))
	claims := DPoPClaims{
		Htm: htm,
		Htu: htu,
		Iat: time.Now().Unix(),
		Jti: uuid.NewString(),
		Ath: base64.RawURLEncoding.EncodeToString(hash[:]),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"htm": claims.Htm,
		"htu": claims.Htu,
		"iat": claims.Iat,
		"jti": claims.Jti,
		"ath": claims.Ath,
	})

	block, _ := pem.Decode([]byte(privPEM))
	if block == nil || block.Type != "EC PRIVATE KEY" {
		d.logger.Error().
			Str("func", "CreateTokenDPopWithAccessToken").
			Msg("failed to decode PEM block containing EC private key")
		return "", errors.New("failed to decode PEM block containing EC private key")
	}

	blockPub, _ := pem.Decode([]byte(pubPEM))
	if blockPub == nil || blockPub.Type != "PUBLIC KEY" {
		d.logger.Error().
			Str("func", "CreateTokenDPopWithAccessToken").
			Msg("failed to decode PEM block containing EC public key")
		return "", errors.New("failed to decode PEM block containing EC public key")
	}

	privateKey, err := x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return "", err
	}

	pubKeyInterface, err := x509.ParsePKIXPublicKey(blockPub.Bytes)
	if err != nil {
		d.logger.Error().Err(err).Msg("failed to parse public key")
		return "", err
	}
	publicKey, ok := pubKeyInterface.(*ecdsa.PublicKey)
	if !ok {
		return "", errors.New("not an ECDSA public key")
	}

	token.Header["typ"] = "dpop+jwt"

	// JWK must contain raw coordinate bytes encoded as Base64URL
	x := publicKey.X.Bytes()
	y := publicKey.Y.Bytes()

	// Ensure padding to the curve's coordinate length (32 bytes for P-256)
	byteLen := (publicKey.Curve.Params().BitSize + 7) / 8
	if len(x) < byteLen {
		padding := make([]byte, byteLen-len(x))
		x = append(padding, x...)
	}
	if len(y) < byteLen {
		padding := make([]byte, byteLen-len(y))
		y = append(padding, y...)
	}

	jwk := map[string]any{
		"kty": "EC",
		"crv": "P-256",
		"x":   base64.RawURLEncoding.EncodeToString(x),
		"y":   base64.RawURLEncoding.EncodeToString(y),
	}

	token.Header["jwk"] = jwk

	signedToken, err := token.SignedString(privateKey)
	if err != nil {
		d.logger.Error().
			Err(err).Send()
		return "", err
	}

	return signedToken, nil
}




func (d *DPoP) IssueAccessTokenwithDpop(token string, 
										method string,
										url string,
										accessToken string,) (error) {
	
	parsed, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
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

		return publicKey, nil
	})
    if err != nil {
        return err
    }
	
	claims := parsed.Claims.(jwt.MapClaims)

    if claims["htm"] != method {
        return errors.New("invalid method")
    }

	if claims["htu"] != url {
        return errors.New("invalid url")
    }

	iat := int64(claims["iat"].(float64))

    if math.Abs(float64(time.Now().Unix()-iat)) > 60 {
        return errors.New("token expired")
    }

	// Cache the jti to prevent replay attacks. In a real implementation, you would store this in a database or cache with an expiration time.
	jti := claims["jti"].(string)
	_ = jti // In a real implementation, you would check if this jti has been seen before.

	hash := sha256.Sum256([]byte(accessToken))
	expectedATH := base64.RawURLEncoding.EncodeToString(hash[:])
	if claims["ath"] != expectedATH {
		return errors.New("invalid access token hash")
	}

	d.logger.Debug().
		Str("func", "ValidateToken").
		Str("method", method).
		Str("url", url).
		Str("iat", time.Unix(iat, 0).String()).
		Str("ath", claims["ath"].(string)).
		Str("expected_ath", expectedATH).
		Str("jti", jti).
		Msg("DPoP token validated successfully")

	return nil
}
