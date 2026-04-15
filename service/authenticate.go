package service

import (
	"fmt"
	"forgeflow-api/config"
	"forgeflow-api/errors"
	"forgeflow-api/repository"
	"time"

	"github.com/MicahParks/keyfunc/v2"
	"github.com/golang-jwt/jwt/v5"
)

type AuthenticateService struct {
	jwks   *keyfunc.JWKS
	config *config.Config
}

type JWTClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthenticateService(userRepo *repository.UserRepository, cfg *config.Config) *AuthenticateService {
	jwksURL := fmt.Sprintf("%s/.well-known/jwks.json", cfg.ClerkFrontendAPI)

	jwks, err := keyfunc.Get(jwksURL, keyfunc.Options{
		RefreshInterval: time.Hour,
	})
	if err != nil {
		panic(fmt.Sprintf("Failed to fetch JWKS: %v", err))
	}

	return &AuthenticateService{jwks: jwks, config: cfg}
}

func (s *AuthenticateService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, s.jwks.Keyfunc,
		jwt.WithIssuer(s.config.ClerkFrontendAPI),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, errors.ErrInvalidToken
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.ErrInvalidToken
}
