package main

import (
  "fmt"
  "github.com/google/uuid"
  "skfw/papaya/pigeon/templates/basicAuth/repository"
)

func main() {

  id := uuid.UUID{}

  fmt.Println(id.String(), repository.EmptyId(id))
}
