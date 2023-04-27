package main

import (
  "fmt"
  "skfw/papaya/pigeon/templates/basic/repository"
  "time"
)

func main() {
  // Example usage
  secret, _ := repository.CreateSecretKey()
  data := map[string]any{
    "sub":  "1234567890",
    "name": "John Doe",
    "iat":  time.Now().Unix(),
    "exp":  time.Now().Add(time.Hour * 24).Unix(),
  }
  tokenString, _ := repository.EncodeJWT(data, secret)
  fmt.Println(tokenString)
  claims, _ := repository.DecodeJWT(tokenString, secret, time.Time{})

  fmt.Println(claims)
}
