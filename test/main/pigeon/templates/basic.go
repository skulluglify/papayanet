package main

import (
  "fmt"
  "skfw/papaya/pigeon/templates/basic"
  "time"
)

func main() {
  // Example usage
  secret, _ := basic.CreateSecretKey()
  data := map[string]any{
    "sub":  "1234567890",
    "name": "John Doe",
    "iat":  time.Now().Unix(),
    "exp":  time.Now().Add(time.Hour * 24).Unix(),
  }
  tokenString, _ := basic.EncodeJWT(data, secret)
  fmt.Println(tokenString)
  claims, valid, _ := basic.DecodeJWT(tokenString, secret)
  if valid {
    fmt.Println(claims)
  }
  fmt.Println(basic.ExpiredJWT(tokenString, time.Now()))
}
