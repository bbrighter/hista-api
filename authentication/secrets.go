package authentication

import (
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"log"
)

var secrets struct {
	PrivateKey string // Private RSA key for JWT
	PublicKey  string // Public RSA key for JWT
}

type Keys struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

var keys *Keys

func newKeys() *Keys {
	dec, _ := pem.Decode([]byte(secrets.PrivateKey))
	if dec == nil {
		log.Fatal("cannot decode private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(dec.Bytes)
	if err != nil {
		log.Fatalf("cannot parse private key: %s", err)
	}
	keys = &Keys{
		privateKey: privateKey,
		publicKey:  &privateKey.PublicKey,
	}
	return keys
}
