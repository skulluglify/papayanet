package main

import (
  "crypto/rand"
  "crypto/rsa"
  "crypto/x509"
  "encoding/pem"
  "fmt"
  "os"
)

func main() {

  // Generate a new private key
  privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
  if err != nil {
    fmt.Println("Error generating private key:", err)
    return
  }

  // Encode the private key in PKCS#8 format
  privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
  privateKeyBlock := pem.Block{
    Type:  "RSA PRIVATE KEY",
    Bytes: privateKeyBytes,
  }
  privateKeyPEM := pem.EncodeToMemory(&privateKeyBlock)

  // Write the private key to a file
  err = os.WriteFile("private.key", privateKeyPEM, 0600)
  if err != nil {
    fmt.Println("Error writing private key file:", err)
    return
  }

  // Extract the public key from the private key
  publicKey := &privateKey.PublicKey

  // Encode the public key in PKIX format
  publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
  if err != nil {
    fmt.Println("Error encoding public key:", err)
    return
  }
  publicKeyBlock := pem.Block{
    Type:  "PUBLIC KEY",
    Bytes: publicKeyBytes,
  }
  publicKeyPEM := pem.EncodeToMemory(&publicKeyBlock)

  // Write the public key to a file
  err = os.WriteFile("public.key", publicKeyPEM, 0644)
  if err != nil {
    fmt.Println("Error writing public key file:", err)
    return
  }

  fmt.Println("Keys generated successfully")
}
