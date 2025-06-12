package crypto

import (
	"crypto/ecdsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
)

func ParseECDSAPrivateKeyFromPemString(prvtPEM []byte) (prvt *ecdsa.PrivateKey, err error) {
	block, _ := pem.Decode(prvtPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	prvt, err = x509.ParseECPrivateKey(block.Bytes)
	if err != nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}
	return
}

func ParseECDSAPublicKeyFromPemString(pubPEM []byte) (*ecdsa.PublicKey, error) {
	block, _ := pem.Decode(pubPEM)
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// Type switch is needed...
	switch pub := pub.(type) {
	case *ecdsa.PublicKey:
		return pub, nil
	default:
		return nil, errors.New("unsupported key type")
	}
}

// A type assertion provides access to an interface value's underlying concrete value. t := i.(T) This statement asserts that the interface value i holds the concrete type T and assigns the underlying T value to the variable t .
