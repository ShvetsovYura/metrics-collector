package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func EncryptData(msg []byte, pubKeyPath string) ([]byte, error) {
	publicKeyBytes, err := os.ReadFile(pubKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error on read public key file %w", err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyBytes)
	publicKey, err := x509.ParsePKIXPublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error on parsed public key %w", err)
	}
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey.(*rsa.PublicKey), msg)
	if err != nil {
		return nil, fmt.Errorf("error on encrypt message %w", err)
	}

	return cipherText, nil
}

func DecryptData(cipherMsg []byte, privateKeyPath string) ([]byte, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error on read private key file %w", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error on parse private key %w", err)
	}
	msg, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherMsg)
	if err != nil {
		return nil, fmt.Errorf("error on decrypt message %w", err)
	}

	return msg, nil

}
