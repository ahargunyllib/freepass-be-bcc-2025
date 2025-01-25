package jwt

import (
	"time"

	"github.com/ahargunyllib/freepass-be-bcc-2025/internal/infra/env"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type CustomJwtInterface interface {
	Create(userID uuid.UUID, role int16) (string, error)
	Decode(tokenString string, claims *Claims) error
}

type Claims struct {
	jwt.RegisteredClaims
	UserID uuid.UUID
	Role   int16
}

type CustomJwtStruct struct {
	SecretKey   string
	ExpiredTime time.Duration
}

var Jwt = getJwt()

func getJwt() CustomJwtInterface {
	return &CustomJwtStruct{
		SecretKey:   env.AppEnv.JwtSecretKey,
		ExpiredTime: env.AppEnv.JwtExpTime,
	}
}

func (j *CustomJwtStruct) Create(userID uuid.UUID, role int16) (string, error) {
	claims := Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "bcc-conference",
			Subject:   userID.String(),
			Audience:  jwt.ClaimStrings{"bcc-conference"},
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.ExpiredTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		UserID: userID,
		Role:   role,
	}

	unsignedJWT := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedJWT, err := unsignedJWT.SignedString([]byte(j.SecretKey))
	if err != nil {
		return "", err
	}

	return signedJWT, nil
}

func (j *CustomJwtStruct) Decode(tokenString string, claims *Claims) error {
	token, err := jwt.ParseWithClaims(tokenString, claims, func(_ *jwt.Token) (any, error) {
		return []byte(j.SecretKey), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return jwt.ErrSignatureInvalid
	}

	return nil
}
