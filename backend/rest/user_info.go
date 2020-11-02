package rest

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

type UserInfo struct {
	ID string
}

func GetUserInfo(request *http.Request) (UserInfo, error) {
	token, err := GetAuthTokenFromHeader(request)
	if err != nil {
		return UserInfo{}, err
	}

	userInfo, err := getUserInfoFromToken(token)
	if err != nil {
		return UserInfo{}, err
	}

	return userInfo, nil
}

func GetAuthTokenFromHeader(r *http.Request) (string, error) {
	rawAccessToken := r.Header.Get("Authorization")

	parts := strings.Split(rawAccessToken, " ")
	if len(parts) != 2 {
		return "", errors.New("wrong authorization header")
	}

	return parts[1], nil
}

func getUserInfoFromToken(token string) (UserInfo, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return UserInfo{}, errors.New("wrong jwt token")
	}

	decodedPayload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return UserInfo{}, err
	}

	dec := json.NewDecoder(bytes.NewReader(decodedPayload))
	var payload payload
	if err := dec.Decode(&payload); err != nil {
		return UserInfo{}, fmt.Errorf("unable to read key %s", err)
	}

	return UserInfo{ID: payload.Sub}, nil
}

type payload struct {
	Exp            int      `json:"exp"`
	Iat            int      `json:"iat"`
	Auth_time      int      `json:"auth_time"`
	Jti            string   `json:"jti"`
	Iss            string   `json:"iss"`
	Aud            string   `json:"aud"`
	Sub            string   `json:"sub"`
	Typ            string   `json:"typ"`
	Azp            string   `json:"azp"`
	Nonce          string   `json:"nonce"`
	SessionState   string   `json:"session_state"`
	Acr            string   `json:"acr"`
	AllowedOrigins []string `json:"allowed-origins"`
	//RealmAccess       []string `json:"realm_access"`
	//ResourceAccess    []string `json:"resource_access"`
	Scope             string `json:"scope"`
	EmailVerified     bool   `json:"email_verified"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
	GivenName         string `json:"given_name"`
	FamilyName        string `json:"family_name"`
	Mail              string `json:"email"`
}

//type RealmAccess struct {
//	Roles []string `json:"roles"`
//}
//
//type ResourceAccess struct {
//	Roles []string `json:"roles"`
//}
