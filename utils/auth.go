package utils

import (
    "errors"
    "time"

    "github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("697FF549C89ED7CEF8D8A97A392A37D47BF5C795B482E89948444CA832") 

type Claims struct {
    CustomerID string `json:"customer_id"`
    jwt.StandardClaims
}

func GenerateJWT(customerID string) (string, error) {
    expirationTime := time.Now().Add(1 * time.Hour)
    claims := &Claims{
        CustomerID: customerID,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: expirationTime.Unix(),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ValidateJWT(tokenStr string) (string, error) {
    claims := &Claims{}
    token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return "", errors.New("invalid signature")
        }
        return "", err
    }
    if !token.Valid {
        return "", errors.New("invalid token")
    }
    return claims.CustomerID, nil
}
