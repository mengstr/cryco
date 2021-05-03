package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
)

const (
	keylen = 16
)

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func main() {
	var err error
	key := make([]byte, keylen)

	genKey := flag.Bool("gen", false, "Generate key")
	keyName := flag.String("key", "", "Use env '<name>' instead of 'crycokey' as key")
	flag.Parse()
	plaintext := flag.Arg(0)

	if *genKey {
		_, err = rand.Read(key)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't generate random key: %s\n", err)
			os.Exit(1)
		}
		fmt.Println(base64.URLEncoding.EncodeToString(key))
		os.Exit(0)
	}

	s := os.Getenv("crycokey")
	if s != "" {
		b, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't decode crycokey from env: %s\n", err)
			os.Exit(1)
		}
		if len(b) != keylen {
			fmt.Fprintf(os.Stderr, "Decoded crycokey is not 16 bytes\n")
			os.Exit(1)
		}
		key = b
	}

	s = os.Getenv(*keyName)
	if *keyName != "" && s == "" {
		fmt.Fprintf(os.Stderr, "Env '%s' dosen't exist or is empty\n", *keyName)
		os.Exit(1)
	}
	if s != "" {
		b, err := base64.URLEncoding.DecodeString(s)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Can't decode key from env '%s': %s\n", s, err)
			os.Exit(1)
		}
		if len(b) != keylen {
			fmt.Fprintf(os.Stderr, "Decoded env '%s' is not 16 bytes\n", *keyName)
			os.Exit(1)
		}
		key = b
	}

	if allZero(key) {
		fmt.Fprintf(os.Stderr, "No key found\n")
		os.Exit(1)
	}

	if plaintext == "" {
		fmt.Fprintf(os.Stderr, "No plaintext specified\n")
		os.Exit(1)
	}

	cipherBlock, err := aes.NewCipher(key)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing cipher(1)\n")
		os.Exit(1)
	}

	aead, err := cipher.NewGCM(cipherBlock)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing cipher(2)\n")
		os.Exit(1)
	}

	nonce := make([]byte, aead.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error initializing cipher(3)\n")
		os.Exit(1)
	}

	fmt.Println(base64.URLEncoding.EncodeToString(aead.Seal(nonce, nonce, []byte(plaintext), nil)))
}
