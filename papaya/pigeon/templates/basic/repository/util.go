package repository

import (
  "crypto/rand"
  "encoding/base64"
  "errors"
  "github.com/golang-jwt/jwe"
  "github.com/golang-jwt/jwt/v5"
  "github.com/valyala/fasthttp"
  "golang.org/x/crypto/bcrypt"
  "strings"
  "time"
)

// generated from chatGPT, more fixed and readability for utilities

// function to hash a password using bcrypt

func HashPassword(password string) (string, error) {

  hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

  if err != nil {

    return "", err
  }

  return string(hashedPassword), nil
}

// function to compare a plaintext password with a hashed password

func CheckPasswordHash(password string, hash string) bool {

  err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))

  return err == nil
}

// CreateSecretKey generates a new random secret key
func CreateSecretKey() (string, error) {

  key := make([]byte, 32)

  _, err := rand.Read(key)

  if err != nil {

    return "", err
  }

  return base64.URLEncoding.EncodeToString(key), nil
}

// EncodeJWT encodes a map of claims into a JWT token string
func EncodeJWT(data map[string]any, secret string) (string, error) {

  claims := jwt.MapClaims(data)

  token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

  return token.SignedString([]byte(secret))
}

// DecodeJWT decodes a JWT token string and returns the claims if successful
func DecodeJWT(token string, secret string, expirationTime time.Time) (jwt.MapClaims, error) {

  var err error
  var ok bool
  var exp float64
  var obj *jwt.Token
  var claims jwt.MapClaims

  claims = jwt.MapClaims{}

  obj, err = jwt.Parse(token, func(token *jwt.Token) (any, error) {
    return []byte(secret), nil
  })

  if err != nil {

    // bypass signature validation
    if err != jwt.ErrSignatureInvalid {

      return nil, errors.New("broken token")
    }
  }

  if claims, ok = obj.Claims.(jwt.MapClaims); ok {

    if !obj.Valid {

      // catch token if not validation by signature
      return claims, errors.New("invalid token signature")
    }

    if exp, ok = claims["exp"].(float64); ok {

      if !time.Unix(int64(exp), 0).Before(expirationTime) {

        return nil, errors.New("expired token")
      }
    }

    return claims, nil
  }

  return nil, errors.New("can't cast token as MapClaims")
}

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

func RequestAuth(req *fasthttp.Request) string {

  var noop string

  auth := string(req.Header.Peek("Authorization"))
  if auth != "" {

    if token, found := strings.CutPrefix(auth, "Bearer "); found {

      return token
    }
  }

  return noop
}

func RequestJWT(req *fasthttp.Request, secret string, expirationTime time.Time) (jwt.MapClaims, error) {

  token := RequestAuth(req)

  if token != "" {

    obj, err := DecodeJWT(token, secret, expirationTime)
    if err != nil {

      return nil, err
    }

    return obj, nil
  }

  return nil, errors.New("no implemented authentication")
}

func RequestJWE(req *fasthttp.Request, privKey []byte, secret string, expirationTime time.Time) (jwt.MapClaims, error) {

  token := RequestAuth(req)

  if token != "" {

    decrypted, err := DecryptJWT(token, privKey)
    if err != nil {

      return nil, err
    }

    obj, err := DecodeJWT(decrypted, secret, expirationTime)
    if err != nil {

      return nil, err
    }

    return obj, nil
  }

  return nil, errors.New("no implemented authentication")
}
