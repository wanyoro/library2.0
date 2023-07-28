package middleware

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"lib2.0/api/responses"
)

// func //SetContentTypeMiddleware sets middleware to json
func SetContentTypeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "Application/json")
		next.ServeHTTP(w, r)
	})
}

// func AuthJwtVerify verifies jwt token
func AuthJwtVerify(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var resp = map[string]interface{}{"status": "failed", "message": "missing Auth token"}

		var bearer = r.Header.Get("Authorization")
		bearer = strings.TrimSpace(bearer)

		if bearer == "" {
			responses.JSON(w, http.StatusForbidden, resp)
			return
		}

		token, err := jwt.Parse(bearer, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("API_KEY")), nil
		})
		if err != nil {
			resp["status"] = "failed"
			resp["message"] = "Invalid Token please login"
			responses.JSON(w, http.StatusForbidden, resp)
			fmt.Printf("%s", err)
			return
		}
		claims, _ := token.Claims.(jwt.MapClaims)

		ctx := context.WithValue(r.Context(), "studentID", claims["studentID"])
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
