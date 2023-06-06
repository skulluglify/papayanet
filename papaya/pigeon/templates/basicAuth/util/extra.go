package util

import (
  "bytes"
  "crypto/rand"
  "encoding/base64"
  "encoding/hex"
  "encoding/json"
  "errors"
  "golang.org/x/crypto/sha3"
  "io"
  "net"
  "net/url"
  "skfw/papaya/ant/bpack"
  "skfw/papaya/bunny/swag"
  "skfw/papaya/koala/pp"
  "skfw/papaya/pigeon/templates/basicAuth/models"
  "strings"
  "time"

  "github.com/golang-jwt/jwe"
  "github.com/golang-jwt/jwt/v5"
  "github.com/google/uuid"
  "github.com/valyala/fasthttp"
  "golang.org/x/crypto/bcrypt"
)

// generated from chatGPT, more fixed and readability for utilities

// function to hash a password using bcrypt

func EmptyAsterisk(t string) string {

  if t != "*" {

    return t
  }

  return ""
}

func EmptyIds(id string) bool {

  for _, c := range []byte(id) {

    if c != 48 { // 00 11 00 00

      return false
    }
  }

  return true
}

func EmptyId(id []byte) bool {

  return EmptyIds(hex.EncodeToString(id))
}

func EmptyIdx(id uuid.UUID) bool {

  var err error
  var idx []byte

  idx, err = id.MarshalBinary()

  if err != nil {

    return false
  }

  return EmptyId(idx)
}

func Id(id []byte) uuid.UUID {

  var err error
  var idx uuid.UUID

  idx, err = uuid.FromBytes(id)
  if err != nil {

    return uuid.UUID{}
  }

  return idx
}

func Idx(id uuid.UUID) string {

  var err error
  var data []byte

  data, err = id.MarshalBinary()
  if err != nil {

    return "00000000000000000000000000000000"
  }

  return hex.EncodeToString(data)
}

func Ids(id string) uuid.UUID {

  var err error
  var data []byte

  // remove all pad
  id = strings.ReplaceAll(id, "-", "")

  data, err = hex.DecodeString(id)
  if err != nil {

    return uuid.UUID{}
  }

  return Id(data)
}

func URLValid(rawURL string) bool {

  if rawURL != "" {

    if _, err := url.Parse(rawURL); err != nil {

      return false
    }

    return true
  }

  return false
}

func GetClientIPFromXForwardedFor(ctx *swag.SwagContext) string {

  XForwardedFor := ctx.Get("X-Forwarded-For")
  tokens := strings.Split(XForwardedFor, ",")
  n := len(tokens)

  if n > 0 {

    for i := 0; i < n; i++ {
      j := n - i - 1

      token := strings.Trim(tokens[j], " ")

      if URLValid(token) {

        return token
      }
    }
  }

  return ""
}

func GetClientIP(ctx *swag.SwagContext) string {

  ClientIPFromRealIP := ctx.Get("X-Real-IP")
  ClientIPFromXForwardedFor := GetClientIPFromXForwardedFor(ctx)
  ClientIP := ctx.IP()

  return pp.Qstr(ClientIPFromRealIP, ClientIPFromXForwardedFor, ClientIP)
}

func CheckIP(sessionClientIP string, currentClientIP string) bool {

  // random ipv4
  // issue about development in emulator

  // insensitive case check, 'c' class
  sIP := net.ParseIP(sessionClientIP)
  cIP := net.ParseIP(currentClientIP)

  // same ipv4
  if sIP.To4() != nil && cIP.To4() != nil {

    sIP = sIP.Mask(net.CIDRMask(32, 32))
    cIP = cIP.Mask(net.CIDRMask(32, 32))

    return sIP.Equal(cIP)
  }

  // iam don't know what have done yet, same ipv6
  if sIP.To16() != nil && cIP.To16() != nil {

    sIP = sIP.Mask(net.CIDRMask(128, 128))
    cIP = cIP.Mask(net.CIDRMask(128, 128))

    return sIP.Equal(cIP)
  }

  return false
}

func DeviceRecognition(ctx *swag.SwagContext, session *models.SessionModel) bool {

  ClientIP := GetClientIP(ctx)
  UserAgent := strings.Trim(ctx.Get("User-Agent"), " ")

  if strings.HasPrefix(session.UserAgent, "Dart/3") ||
    strings.HasPrefix(session.UserAgent, "Java/1") {
    // issue test on emulator 'Dart/3.0 (dart:io)'
    // randomize Client IP on emulator
    return session.UserAgent == UserAgent
  }

  // Aggressive Checker
  return CheckIP(session.ClientIP, ClientIP) &&
    session.UserAgent == ctx.Get("User-Agent")
}

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
func DecodeJWT(token string, secret string) (jwt.MapClaims, error) {

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

    if err == jwt.ErrTokenExpired {

      return nil, errors.New("expired token")
    }

    // bypass signature validation
    if err != jwt.ErrSignatureInvalid {

      return nil, errors.New("broken token")
    }
  }

  currentTime := time.Now().UTC()

  if claims, ok = obj.Claims.(jwt.MapClaims); ok {

    if !obj.Valid {

      // catch token if not validation by signature
      return claims, errors.New("invalid token signature")
    }

    if exp, ok = claims["exp"].(float64); ok {

      expiredTime := time.UnixMilli(int64(exp))

      if currentTime.After(expiredTime) {

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

    // try with a lower case
    if len(auth) > 7 {

      bearer := strings.ToLower(auth[:6])
      if bearer == "bearer" {

        return auth[7:]
      }
    }

    // bypass
    return auth
  }

  return noop
}

func RequestJWT(req *fasthttp.Request, secret string) (jwt.MapClaims, error) {

  token := RequestAuth(req)

  if token != "" {

    obj, err := DecodeJWT(token, secret)
    if err != nil {

      return nil, err
    }

    return obj, nil
  }

  return nil, errors.New("no implemented authentication")
}

func RequestJWE(req *fasthttp.Request, privKey []byte, secret string) (jwt.MapClaims, error) {

  token := RequestAuth(req)

  if token != "" {

    decrypted, err := DecryptJWT(token, privKey)
    if err != nil {

      return nil, err
    }

    obj, err := DecodeJWT(decrypted, secret)
    if err != nil {

      return nil, err
    }

    return obj, nil
  }

  return nil, errors.New("no implemented authentication")
}

func GetTLDs() []string {

  var data []string
  data = make([]string, 0)

  if packet := bpack.OpenPacket("/data/kornet/tlds.json"); packet != nil {

    if err := json.Unmarshal(packet.Data, &data); err != nil {

      return nil
    }
  }

  return data
}

func TLDChecker(tlds []string, address string) bool {

  if tlds != nil {

    tokens := strings.Split(address, ".")
    n := len(tokens)

    if n > 1 {

      suffix := tokens[n-1]

      for _, tld := range tlds {

        if strings.ToUpper(tld) == strings.ToUpper(suffix) {

          return true
        }
      }
    }
  }

  return false
}

func HashSHA3(input string) string {

  hash := sha3.New256()
  io.WriteString(hash, input)

  return hex.EncodeToString(hash.Sum(nil))
}

func HashCompareSHA3(input string, hash string) bool {

  var err error
  var a, b []byte
  if a, err = hex.DecodeString(HashSHA3(input)); err != nil {

    return false
  }
  if b, err = hex.DecodeString(hash); err != nil {

    return false
  }

  return bytes.Equal(a, b)
}
