package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func main() {
	key, err := Generate25519Key()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	WriteKeyToFile(key)
	os.Exit(0)
}

func WriteKeyToFile(key *Key) error {
	write := func(name string, block *pem.Block) error {
		f, err := os.Create(name)
		if err != nil {
			return fmt.Errorf("failed to create file: %w", err)
		}
		defer f.Close()
		_, err = f.Write(pem.EncodeToMemory(block))
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
		return nil
	}

	buf, err := x509.MarshalPKCS8PrivateKey(key.prv)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	block := &pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: buf,
	}

	err = write("private.ed25519.pem", block)
	if err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	buf, err = x509.MarshalPKIXPublicKey(key.pub)
	if err != nil {
		return fmt.Errorf("failed to marshal public key: %w", err)
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: buf,
	}
	err = write("public.ed25519.pem", block)
	if err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}
	return nil
}

type Key struct {
	pub ed25519.PublicKey
	prv ed25519.PrivateKey
}

func Generate25519Key() (*Key, error) {
	// generate Generate25519Key
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}
	return &Key{public, private}, nil
}
