package base

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/libp2p/go-libp2p/core/crypto"
	"os"
)

// generatePrivateKey creates a new Ed25519 private key and returns it in serialized form.
func GeneratePrivateKey() (crypto.PrivKey, error) {
	priv, _, err := crypto.GenerateEd25519Key(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %v", err)
	}
	return priv, nil
}

// loadPrivateKeyFromEnv loads a private key from an environment variable.
func LoadPrivateKeyFromEnv(passphrase string) (crypto.PrivKey, error) {
	keyPEM := os.Getenv("NODE_PRIVATE_KEY")
	if keyPEM == "" {
		return nil, fmt.Errorf("NODE_PRIVATE_KEY environment variable not set")
	}

	block, _ := pem.Decode([]byte(keyPEM))
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	// Check if the key is encrypted
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		decryptedBytes, err := x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt key: %v", err)
		}
		block.Bytes = decryptedBytes
	}
	privKey, err := crypto.UnmarshalPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %v", err)
	}
	return privKey, nil
}

// loadEncryptedPrivateKeyFromFile loads an encrypted private key from a file.
func LoadEncryptedPrivateKeyFromFile(filePath, passphrase string) (crypto.PrivKey, error) {
	encryptedData, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read key file: %v", err)
	}

	block, _ := pem.Decode(encryptedData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	// Check if the key is encrypted
	if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
		decryptedBytes, err := x509.DecryptPEMBlock(block, []byte(passphrase))
		if err != nil {
			return nil, fmt.Errorf("failed to decrypt key: %v", err)
		}
		block.Bytes = decryptedBytes
	}

	privKey, err := crypto.UnmarshalPrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal private key: %v", err)
	}
	return privKey, nil
}

// savePrivateKeyToEncryptedFile saves a private key to an encrypted PEM file.
func SavePrivateKeyToEncryptedFile(privKey crypto.PrivKey, filePath, passphrase string) error {
	keyBytes, err := crypto.MarshalPrivateKey(privKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %v", err)
	}

	// Encrypt the key
	encryptedBlock, err := x509.EncryptPEMBlock(rand.Reader, "ENCRYPTED PRIVATE KEY", keyBytes, []byte(passphrase), x509.PEMCipherAES256)
	if err != nil {
		return fmt.Errorf("failed to encrypt private key: %v", err)
	}

	pemData := pem.EncodeToMemory(encryptedBlock)
	return os.WriteFile(filePath, pemData, 0600) // Restrict permissions
}
