package rest

import (
	"net/http"
)

const adminRole = "admin"
const projectManagerRole = "project-manager"

type AuthService struct {
	authenticator Authenticator
}

type Authenticator interface {
	verify(token string) error
}

func (a *AuthService) HttpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := GetAuthTokenFromHeader(r)
		if err != nil {
			SendAuthorizationErrorJSON(w, r, err)
			return
		}

		err = a.authenticator.verify(token)
		if err != nil {
			SendAuthorizationErrorJSON(w, r, err)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (a *AuthService) HasRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := GetAuthTokenFromHeader(r)
			if err != nil {
				SendAuthorizationErrorJSON(w, r, err)
				return
			}

			userInfo, err := getUserInfoFromToken(token)
			if err != nil || !userInfo.hasRole(role) {
				SendAuthorizationErrorJSON(w, r, err)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
