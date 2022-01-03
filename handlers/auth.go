package handlers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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
    a.l.Println("handle LOGIN request")
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
        http.Error(rw, "no basic auth", http.StatusUnauthorized)
        return
    }

    // are the details correct?
    if username != os.Getenv("USERNAME") || password != os.Getenv("PASSWORD") {
        http.Error(rw, "incorrect username or password", http.StatusForbidden)
        return
    }

    adminJWT, _, err := a.generateJWT()
    if err != nil {
        http.Error(rw, "error generating jwt", http.StatusInternalServerError)
    }

    rw.Header().Add("Access_token", adminJWT)
}

func (a Auth) MiddlewareAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
        c, err := r.Cookie("token")
        if err != nil {
            if err == http.ErrNoCookie {
                http.Error(rw, "No JWT cookie", http.StatusUnauthorized)
                return
            }
            http.Error(rw, "error getting cookie", http.StatusBadRequest)
        }

        tknStr := c.Value

        token, err := jwt.Parse(tknStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }

            return []byte(os.Getenv("HMACSECRET")), nil
        })

        if err != nil {
            http.Error(rw, fmt.Sprintf("Error with jwt: %v", err), http.StatusForbidden)
            return
        }

        if token.Valid {
            next.ServeHTTP(rw, r)
        } else {
            http.Error(rw, "JWT invalid", http.StatusForbidden)
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
        a.l.Printf("error generating tokenstring: %s", err)
        return "", time.Now(), err
    }

    return tokenString, expTime, nil
}
