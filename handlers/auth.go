package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Jedsonofnel/hubdc-api/data"
	jwt "github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"
)

type Auth struct {
    l *log.Logger
}

func NewAuth(l *log.Logger) *Auth {
    return &Auth{l}
}

func (a Auth) Login(rw http.ResponseWriter, r *http.Request) {
    a.l.Println("Handling LOGIN request")

    username, password, ok := r.BasicAuth()

    // godotenv for env file
    if os.Getenv("APP_ENV") != "production" {
        err := godotenv.Load()
        if err != nil {
            a.l.Fatal("Error loading .env file")
        }
    }

    // is auth header in request?
    if !ok {
        a.l.Printf("Error handling LOGIN request: No basic auth present")
        a.JSONError(
            rw,
            data.NewJE("No basic auth set"),
            http.StatusUnauthorized,
        )
        return
    }

    // are the details correct?
    if username != os.Getenv("USERNAME") || password != os.Getenv("PASSWORD") {
        a.l.Printf("Error handling LOGIN request: Incorrect username or password")
        a.JSONError(
            rw,
            data.NewJE("Incorrect username or password"),
            http.StatusForbidden,
        )
        return
    }

    adminJWT, _, err := a.generateJWT()
    if err != nil {
        a.l.Printf("Error handling LOGIN request: %v", err)
        a.JSONError(
            rw,
            data.NewJE("Error generating JWT"),
            http.StatusInternalServerError,
        )
        return
    }

    rw.Header().Add("Access_token", adminJWT)
    rw.WriteHeader(http.StatusOK)
}

func (a Auth) MiddlewareAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        tknStr := r.Header.Get("Authorization")
        if tknStr == "" {
            a.l.Printf("Error authorizing request: No 'Authorization' header set")
            a.JSONError(
                rw,
                data.NewJE("No 'Authorization' header set"),
                http.StatusUnauthorized,
            )
            return
        }

        token, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            return []byte(os.Getenv("HMACSECRET")), nil
        })

        if err != nil {
            a.l.Printf("Error authorizing request: %v", err)
            a.JSONError(
                rw,
                data.NewJE(fmt.Sprintf("Error parsing JWT: %v", err)),
                http.StatusForbidden,
            )
            return
        }

        if token.Valid {
            next.ServeHTTP(rw, r)
        } else {
            a.l.Printf("Error authorizing request: JWT invalid")
            a.JSONError(
                rw,
                data.NewJE("JWT invalid"),
                http.StatusForbidden,
            )
            return
        }
    })
}

func (a Auth) generateJWT() (string, time.Time, error) {
    expTime := time.Now().Add(time.Hour * 1)
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user": "HCM",
        "exp": expTime.Unix(),
    })

    tokenString, err := token.SignedString([]byte(os.Getenv("HMACSECRET")))
    if err != nil {
        a.l.Printf("Error generating tokenstring: %v", err)
        return "", time.Now(), err
    }

    return tokenString, expTime, nil
}
