package main

import (
	"encoding/pem"
	"fmt"
	"github.com/flinkcoin/mono/libs/shared/pkg/base"
	"github.com/libp2p/go-libp2p/core/crypto"
	"os"
)

func main() {
	env, err2 := base.LoadPrivateKeyFromEnv("my-secret-passphrase")
	if err2 != nil {
		panic(err2)
	}

	base.Log.Info("Loaded private key from env:", "type", env.Type().String())
	fmt.Println("Loaded private key from file:", env.Type())
	// Step 4: Load from encrypted file (more secure option)
	fileKey, err := base.LoadEncryptedPrivateKeyFromFile("node_key.pem", "my-secret-passphrase")
	if err != nil {
		panic(err)
	}

	fmt.Println("Loaded private key from file:", fileKey.Type())

}
func main1() {

	// Step 1: Generate a new private key (for demo purposes)
	privKey, err := base.GeneratePrivateKey()
	if err != nil {
		panic(err)
	}
	base.Log.Info("Generated new Ed25519 private key")

	// Step 2: Save it to an encrypted file (optional, for persistent storage)
	err = base.SavePrivateKeyToEncryptedFile(privKey, "node_key.pem", "my-secret-passphrase")
	if err != nil {
		panic(err)
	}
	base.Log.Info("Saved encrypted private key to node_key.pem")

	// Step 3: Load from environment variable (simpler option)
	// Simulate setting the env var (in real use, set it externally)
	keyBytes, _ := crypto.MarshalPrivateKey(privKey)
	pemData := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: keyBytes})
	os.Setenv("NODE_PRIVATE_KEY", string(pemData))

	envKey, err := base.LoadPrivateKeyFromEnv("my-secret-passphrase")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded private key from env:", envKey.Type())

	// Step 4: Load from encrypted file (more secure option)
	fileKey, err := base.LoadEncryptedPrivateKeyFromFile("node_key.pem", "my-secret-passphrase")
	if err != nil {
		panic(err)
	}
	fmt.Println("Loaded private key from file:", fileKey.Type())

	// Example usage: Get public key or sign something
	pubKey := privKey.GetPublic()
	fmt.Println("Public key type:", pubKey.Type())
}
