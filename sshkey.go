package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	gossh "golang.org/x/crypto/ssh"
)

func generateRSA4096Key(filename string) error {
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyPEM := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	return nil
}

func generateECDSAKey(filename string) error {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return fmt.Errorf("failed to generate ECDSA private key: %v", err)
	}

	privateKeyFile, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %v", err)
	}
	defer privateKeyFile.Close()

	privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal ECDSA private key: %v", err)
	}

	privateKeyPEM := &pem.Block{
		Type:  "EC PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	if err := pem.Encode(privateKeyFile, privateKeyPEM); err != nil {
		return fmt.Errorf("failed to write private key to file: %v", err)
	}

	return nil
}

func loadOrGenerateSSHKey(filename string) (gossh.Signer, error) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Println("Private key file not found, generating a new one...")
		if err := generateECDSAKey(filename); err != nil {
			return nil, err
		}
	}

	privateBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to load private key: %v", err)
	}

	private, err := gossh.ParsePrivateKey(privateBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %v", err)
	}

	return private, nil
}
