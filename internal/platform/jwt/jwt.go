package jwt

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/reverie-jp/piamap/internal/platform/ulid"
)

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrInvalidTokenType = errors.New("invalid token type")
	ErrExpiredToken     = errors.New("expired token")
	ErrInvalidSignature = errors.New("invalid signature")
)

type TokenType string

const (
	TokenTypeAccess  TokenType = "access"
	TokenTypeRefresh TokenType = "refresh"
)

type Claims struct {
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

type Manager struct {
	secretKey         string
	accessExpiration  time.Duration
	refreshExpiration time.Duration
}

func NewManager(secretKey string, accessExpiration, refreshExpiration time.Duration) *Manager {
	return &Manager{
		secretKey:         secretKey,
		accessExpiration:  accessExpiration,
		refreshExpiration: refreshExpiration,
	}
}

func (m *Manager) GenerateAccessToken(userID ulid.ULID) (string, error) {
	claims := &Claims{
		TokenType: TokenTypeAccess,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(m.accessExpiration)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(m.secretKey))
}

func (m *Manager) GenerateRefreshToken(userID ulid.ULID) (string, time.Time, error) {
	expireTime := time.Now().Add(m.refreshExpiration)
	claims := &Claims{
		TokenType: TokenTypeRefresh,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expireTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expireTime, nil
}

func (m *Manager) VerifyToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidSignature
		}
		return []byte(m.secretKey), nil
	}, jwt.WithExpirationRequired())

	if err != nil {
		var claims *Claims
		if token != nil {
			if t, ok := token.Claims.(*Claims); ok {
				claims = t
			}
		}
		if claims != nil && claims.ExpiresAt != nil && time.Now().After(claims.ExpiresAt.Time) {
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, ErrInvalidToken
	}

	switch claims.TokenType {
	case TokenTypeAccess, TokenTypeRefresh:
	default:
		return nil, ErrInvalidTokenType
	}

	return claims, nil
}
