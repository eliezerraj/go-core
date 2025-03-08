package cert

import (
	"crypto/rsa"
    "encoding/pem"
	"crypto/x509"
	"errors"

	"github.com/rs/zerolog/log"
)

var childLogger = log.With().Str("go-core", "cert").Logger()

type CertCore struct {
}

func (c *CertCore) ParsePemToRSAPriv(private_key *string) (*rsa.PrivateKey, error){
	childLogger.Debug().Msg("ParsePemToRSAPriv")
	childLogger.Debug().Interface("private_key :",private_key).Msg("")

	block, _ := pem.Decode([]byte(*private_key))
	if block == nil || block.Type != "PRIVATE KEY" {
		//childLogger.Error().Err(err).Msg("erro PRIVATE KEY Decode")
		return nil, errors.New("erro PRIVATE KEY Decode")
	}

	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro ParsePKCS8PrivateKey")
		return nil, err
	}

	key_rsa := privateKey.(*rsa.PrivateKey)

	return key_rsa, nil
}

func (c *CertCore) ParsePemToRSAPub(public_key *string) (*rsa.PublicKey, error){
	childLogger.Debug().Msg("ParsePemToRSAPub")
	childLogger.Debug().Interface("public_key :",public_key).Msg("")

	block, _ := pem.Decode([]byte(*public_key))
	if block == nil || block.Type != "PUBLIC KEY" {
		childLogger.Error().Err(errors.New("erro PUBLIC KEY Decode")).Msg("erro PUBLIC KEY Decode")
		return nil, errors.New("erro PUBLIC KEY Decode")
	}

	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		childLogger.Error().Err(err).Msg("erro ParsePKCS8PrivateKey")
		return nil, err
	}

	key_rsa := pubInterface.(*rsa.PublicKey)

	return key_rsa, nil
}