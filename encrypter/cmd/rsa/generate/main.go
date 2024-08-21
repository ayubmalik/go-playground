package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func main() {
	generateRSAKey()
}

func generateRSAKey() error {
	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	}

	privFile, err := os.Create("private.pem")
	if err != nil {
		return err
	}
	defer privFile.Close()

	err = pem.Encode(privFile, block)
	if err != nil {
		return err
	}

	block = &pem.Block{
		Type:  "RSA PUBLIC",
		Bytes: x509.MarshalPKCS1PublicKey(&key.PublicKey),
	}

	pubFile, err := os.Create("public.pem")
	if err != nil {
		return err
	}
	defer pubFile.Close()

	err = pem.Encode(pubFile, block)
	if err != nil {
		return err
	}

	return nil
}
