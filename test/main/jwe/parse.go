package main

import (
  "errors"
  "fmt"
  "github.com/golang-jwt/jwe"
  "os"
)

func EncryptJWT(token string, pubKey []byte) (string, error) {

  var noop string

  pk, err := jwe.ParseRSAPublicKeyFromPEM(pubKey)
  if err != nil {

    return noop, errors.New("can't parse RSA pub key")
  }

  obj, err := jwe.NewJWE(jwe.KeyAlgorithmRSAOAEP, pk, jwe.EncryptionTypeA256GCM, []byte(token))
  if err != nil {

    return noop, errors.New("can't encrypt token")
  }

  serialize, err := obj.CompactSerialize()
  if err != nil {

    return noop, errors.New("can't serialize object token")
  }

  return serialize, nil
}

func DecryptJWT(token string, privKey []byte) (string, error) {

  var noop string

  obj, err := jwe.ParseEncrypted(token)
  if err != nil {

    return noop, errors.New("invalid token")
  }

  pk, err := jwe.ParseRSAPrivateKeyFromPEM(privKey)
  if err != nil {

    return noop, errors.New("can't parse RSA priv key")
  }

  decrypted, err := obj.Decrypt(pk)
  if err != nil {

    return noop, errors.New("can't decrypt token")
  }

  return string(decrypted), nil
}

func test() {

  // decrypt token

  var tokenString string

  tokenString = "{ 'name': 'ahmad asy syafiq' }"

  pubKey, _ := os.ReadFile("public.key")
  token, _ := EncryptJWT(tokenString, pubKey)

  fmt.Println("token", token)

  privKey, _ := os.ReadFile("private.key")
  res, _ := DecryptJWT(token, privKey)

  fmt.Println("res", res)
}
