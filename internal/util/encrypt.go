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
		return nil, fmt.Errorf("error on read public key file %e", err)
	}
	publicKeyBlock, _ := pem.Decode(publicKeyBytes)
	publicKey, err := x509.ParsePKCS1PublicKey(publicKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error on parsed public key %e", err)
	}
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, publicKey, msg)
	if err != nil {
		return nil, fmt.Errorf("error on encrypt message %e", err)
	}

	return cipherText, nil
}

func DecryptData(cipherMsg []byte, privateKeyPath string) ([]byte, error) {
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, fmt.Errorf("error on read private key file %e", err)
	}
	privateKeyBlock, _ := pem.Decode(privateKeyBytes)
	privateKey, err := x509.ParsePKCS1PrivateKey(privateKeyBlock.Bytes)
	if err != nil {
		return nil, fmt.Errorf("error on parse private key %e", err)
	}
	msg, err := rsa.DecryptPKCS1v15(rand.Reader, privateKey, cipherMsg)
	if err != nil {
		return nil, fmt.Errorf("error on decrypt message %e", err)
	}

	return msg, nil

}
