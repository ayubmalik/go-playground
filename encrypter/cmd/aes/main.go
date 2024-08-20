package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

const (
	IV_SIZE     = 16
	BUFFER_SIZE = 4096
)

func main() {
	secret := []byte("helloworld_helloworld_helloworld")

	start := time.Now()
	encrypt("roger.webp", "roger.enc.webp", secret)
	encTime := time.Since(start)

	start = time.Now()
	decrypt("roger.enc.webp", "roger.dec.webp", secret)
	decTime := time.Since(start)

	fmt.Println("encrypt took", encTime)
	fmt.Println("decrypt took", decTime)
}

func encrypt(inFile, outFile string, key []byte) error {
	in, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	iv := make([]byte, IV_SIZE)
	_, err = rand.Read(iv)
	if err != nil {
		return err
	}
	out.Write(iv)

	aes, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aes, iv)
	src := make([]byte, BUFFER_SIZE)
	dst := make([]byte, BUFFER_SIZE)
	for {
		n, err := in.Read(src)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		ctr.XORKeyStream(dst, src[:n])
		out.Write(dst[:n])

		if err == io.EOF {
			break
		}

	}
	return nil
}

func decrypt(inFile, outFile string, key []byte) error {
	in, err := os.Open(inFile)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(outFile)
	if err != nil {
		return err
	}
	defer out.Close()

	iv := make([]byte, IV_SIZE)
	_, err = in.Read(iv)
	if err != nil {
		return err
	}

	aes, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	ctr := cipher.NewCTR(aes, iv)

	src := make([]byte, BUFFER_SIZE)
	dst := make([]byte, BUFFER_SIZE)
	for {
		n, err := in.Read(src)
		if err != nil && !errors.Is(err, io.EOF) {
			return err
		}
		ctr.XORKeyStream(dst, src[:n])
		out.Write(dst[:n])

		if err == io.EOF {
			break
		}

	}

	return nil
}
