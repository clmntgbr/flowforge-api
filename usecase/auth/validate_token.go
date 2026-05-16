package auth

import (
	"context"
	"errors"
	"flowforge-api/domain/repository"
	authdto "flowforge-api/infrastructure/auth"
	"flowforge-api/infrastructure/clerk"

	"github.com/golang-jwt/jwt/v5"
)

type ValidateTokenUseCase struct {
	jwksProvider *clerk.JWKSProvider
	userRepo     *repository.UserRepository
}

func NewValidateTokenUseCase(
	jwksProvider *clerk.JWKSProvider,
	userRepo *repository.UserRepository,
) *ValidateTokenUseCase {
	return &ValidateTokenUseCase{
		jwksProvider: jwksProvider,
		userRepo:     userRepo,
	}
}

func (uc *ValidateTokenUseCase) Execute(ctx context.Context, input authdto.ValidateTokenInput) (*authdto.ValidateTokenOutput, error) {
	token, err := jwt.ParseWithClaims(
		input.Token,
		&authdto.JWTClaims{},
		uc.jwksProvider.GetKeyfunc(),
		jwt.WithIssuer(uc.jwksProvider.GetIssuer()),
		jwt.WithExpirationRequired(),
	)

	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*authdto.JWTClaims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}

	user, err := (*uc.userRepo).GetByClerkID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &authdto.ValidateTokenOutput{
		User:   user,
		Claims: claims,
	}, nil
}
