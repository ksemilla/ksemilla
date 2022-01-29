package middlewares

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"ksemilla/database"
	"ksemilla/graph/model"

	"github.com/golang-jwt/jwt"
)

type key int

const (
	UserCtx key = iota
)

func GetUserCtx() key {
	return UserCtx
}

var db = database.Connect()

func UserContextBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {

		fmt.Println("AUTH", r.Header.Get("Authorization"))
		authToken := r.Header.Get("Authorization")
		userObj := model.User{}
		if len(authToken) == 0 {
			ctx := context.WithValue(r.Context(), UserCtx, nil)
			next.ServeHTTP(rw, r.WithContext(ctx))
		} else {
			token, err := jwt.Parse(authToken, func(token *jwt.Token) (interface{}, error) {
				// Don't forget to validate the alg is what you expect:
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					// return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
					return nil, errors.New("unexpected signing method")
				}

				// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
				return []byte("Jebaited"), nil
			})

			if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
				timeValue := int64(claims["ExpiresAt"].(float64)) - time.Now().Unix()

				if timeValue <= 0 {
					// return "", errors.New("expired token")
					fmt.Println("expired token")
					ctx := context.WithValue(r.Context(), UserCtx, nil)
					next.ServeHTTP(rw, r.WithContext(ctx))
				} else {
					user, err := db.FindOneUser(claims["userId"].(string))
					if err != nil {
						fmt.Println("auth:no user found with id")
						ctx := context.WithValue(r.Context(), UserCtx, nil)
						next.ServeHTTP(rw, r.WithContext(ctx))
					} else {
						// IF EVERYTHING LOOKS FINE
						userObj = *user
						ctx := context.WithValue(r.Context(), UserCtx, userObj)
						next.ServeHTTP(rw, r.WithContext(ctx))
					}
				}
			} else {
				fmt.Println(err)
				// return "", errors.New("token unrecognized")
				ctx := context.WithValue(r.Context(), UserCtx, nil)
				next.ServeHTTP(rw, r.WithContext(ctx))

			}
		}

		// js, _ := json.Marshal(&struct{ Email string }{"test@test.com"})
		// ctx := context.WithValue(r.Context(), "user", js)
		// next.ServeHTTP(rw, r.WithContext(ctx))
	})
}
