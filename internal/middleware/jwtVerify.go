package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/DroidZed/go_lance/internal/config"
	"github.com/DroidZed/go_lance/internal/cryptor"
	"github.com/DroidZed/go_lance/internal/utils"
)

func AccessVerify(next http.Handler) http.Handler {
	env := config.LoadEnv()
	return tokenVerify(env.AccessSecret, next)
}

func RefreshVerify(next http.Handler) http.Handler {
	env := config.LoadEnv()
	return tokenVerify(env.RefreshSecret, next)
}

func tokenVerify(secret string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := config.GetLogger()

		tokenString, err := retrieveTokenFromHeader(r)

		if err != nil {
			log.Errorf("Header corrupted.\n%s", err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		token, err := cryptor.ParseToken(tokenString, secret)
		if err != nil {
			log.Errorf("Unable to parse token.\n%s", err)
			utils.JsonResponse(
				w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		exp, err := cryptor.ExtractExpiryFromClaims(token)
		if err != nil {
			log.Error(err)
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}

		if expired := time.Unix(int64(exp), 0).Before(time.Now()); expired {
			utils.JsonResponse(w,
				http.StatusUnauthorized,
				utils.DtoResponse{Error: "Invalid token."},
			)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func retrieveTokenFromHeader(r *http.Request) (string, error) {
	header := r.Header.Get("Authorization")

	segments := strings.Split(header, " ")
	if len(segments) != 2 || segments[0] != "Bearer" {
		return "", fmt.Errorf("invalid format")
	}

	return segments[1], nil
}
