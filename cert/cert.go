package cert

import (
	"crypto/rsa"
    "encoding/pem"
	"crypto/x509"
	"errors"

	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("component","go-core").Str("package", "cert.cert").Logger()

type CertCore struct {
}

// About convert a key pem string in rsa key
func (c *CertCore) ParsePemToRSAPriv(private_key *string) (*rsa.PrivateKey, error){
	childLogger.Debug().Str("func","ParsePemToRSAPriv").Interface("private_key :",private_key).Send()

	block, _ := pem.Decode([]byte(*private_key))
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, errors.New("erro PRIVATE KEY Decode")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, err
	}

	key_rsa := privateKey.(*rsa.PrivateKey)

	return key_rsa, nil
}

// About convert a key pem string in rsa key
func (c *CertCore) ParsePemToRSAPub(public_key *string) (*rsa.PublicKey, error){
	childLogger.Debug().Msg("ParsePemToRSAPub")
	childLogger.Debug().Interface("public_key :",public_key).Msg("")

	block, _ := pem.Decode([]byte(*public_key))
	if block == nil || block.Type != "PUBLIC KEY" {
		childLogger.Error().Err(errors.New("erro PUBLIC KEY Decode")).Send()
		return nil, errors.New("erro PUBLIC KEY Decode")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Send()
		return nil, err
	}

	key_rsa := pubInterface.(*rsa.PublicKey)

	return key_rsa, nil
}

// About convert a cert pem string in cert x509
func (c *CertCore) ParsePemToCertx509(certX509pem *string) (*x509.Certificate, error) {
    childLogger.Debug().Msg("ParsePemToCertx509")
	childLogger.Debug().Interface("certX509pem :",certX509pem).Msg("")

	block, _ := pem.Decode([]byte(*certX509pem))
	if block == nil || block.Type != "CERTIFICATE" {
		childLogger.Error().Err(errors.New("erro CERT X509 Decode")).Send()
		return nil, errors.New("erro CERT X509 Decode")
	}

	certX509, err := x509.ParseCertificate(block.Bytes)
    if err != nil {
		childLogger.Error().Err(err).Send()
        return nil, err
    }

	return certX509, nil
}