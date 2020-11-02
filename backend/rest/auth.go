package rest

import (
	"bytes"
	"crypto"
	"crypto/rsa"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const certInfoURL = "http://localhost:9020/auth/realms/master/protocol/openid-connect/certs"

type JWKS struct {
	Keys []JWK
}

type JWK struct {
	Alg string
	Kty string
	X5c []string
	N   string
	E   string
	Kid string
	X5t string
}

type AuthService struct {
	publicKey *rsa.PublicKey
}

func CreateAuthService() *AuthService {
	jwks, err := getCerts()
	if err != nil {
		panic("cannot get certificate")
	}

	return &AuthService{publicKey: getPublicKeyFromJWK(jwks.Keys[0])}
}

func (a *AuthService) HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token, err := GetAuthTokenFromHeader(r); err != nil || a.verifyToken(token) != nil {
			SendErrorJSON(w, r, http.StatusForbidden, err, "user should be logged in", ErrInternal)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *AuthService) verifyToken(token string) error {
	return verifySignature(token, a.publicKey)
}

func verifySignature(jwtToken string, publicKey *rsa.PublicKey) error {
	parts := strings.Split(jwtToken, ".")
	if len(parts) != 3 {
		errors.New("wrong jwt token")
	}

	message := []byte(strings.Join(parts[0:2], "."))
	signature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return err
	}

	// Only small messages can be signed directly; thus the hash of a message, rather than the message itself, is signed.
	hasher := crypto.SHA256.New()
	hasher.Write(message)

	err = rsa.VerifyPKCS1v15(publicKey, crypto.SHA256, hasher.Sum(nil), signature)
	return err
}

func getPublicKeyFromJWK(key JWK) *rsa.PublicKey {
	n, _ := base64.RawURLEncoding.DecodeString(key.N)
	e, _ := base64.RawURLEncoding.DecodeString(key.E)
	z := new(big.Int)
	z.SetBytes(n)

	var buffer bytes.Buffer
	buffer.WriteByte(0)
	buffer.Write(e)
	exponent := binary.BigEndian.Uint32(buffer.Bytes())
	return &rsa.PublicKey{N: z, E: int(exponent)}
}

func getCerts() (JWKS, error) {
	client := http.Client{
		Timeout: time.Duration(5 * time.Second),
	}

	request, err := http.NewRequest("GET", certInfoURL, nil)
	if err != nil {
		return JWKS{}, err
	}

	response, err := client.Do(request)
	if err != nil {
		return JWKS{}, err
	}

	defer response.Body.Close()

	dec := json.NewDecoder(response.Body)
	var jwks JWKS
	if err := dec.Decode(&jwks); err != nil {
		return JWKS{}, fmt.Errorf("unable to read key %s", err)
	}
	return jwks, nil
}
