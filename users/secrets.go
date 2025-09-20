package users

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

type Secrets struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

func NewSecrets() *Secrets {
	dec, _ := pem.Decode([]byte(secrets.PrivateKey))
	if dec == nil {
		log.Fatal("cannot decode private key")
	}
	privateKey, err := x509.ParsePKCS1PrivateKey(dec.Bytes)
	if err != nil {
		log.Fatalf("cannot parse private key: %s", err)
	}

	// pubDec, _ := pem.Decode([]byte(secrets.PublicKey))
	// if pubDec == nil {
	// 	log.Fatal("cannot decode private key")
	// }
	// publicKey, err := x509.ParsePKCS1PublicKey(pubDec.Bytes)
	// if err != nil {
	// 	log.Fatalf("cannot parse public key: %s", err)
	// }
	public := privateKey.PublicKey

	return &Secrets{
		privateKey: privateKey,
		publicKey:  &public,
	}
}
