package crypto

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

// Parse private key from private pem string.
func ParseRSAPrivateKeyFromPemString(prvtKeyPem []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(prvtKeyPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	prvt, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return prvt, nil
}

// Parse public key from pem string.
func ParseRSAPublicKeyFromPemString(pubKeyPem []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pubKeyPem)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return pub, nil
}
