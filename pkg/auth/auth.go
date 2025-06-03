package auth

import (
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/IzePhanthakarn/go-basic-shop/config"
	"github.com/IzePhanthakarn/go-basic-shop/modules/users"
	"github.com/golang-jwt/jwt/v5"
)

type TokenType string

const (
	Access  TokenType = "access"
	Refresh TokenType = "refresh"
	Admin   TokenType = "admin"
	ApiKey  TokenType = "apikey"
)

type basicAuth struct {
	mapClaims *basicMapClaims
	cfg       config.IJwtConfig
}

type basicAdmin struct {
	*basicAuth
}

type basicApiKey struct {
	*basicAuth
}

type basicMapClaims struct {
	Claims *users.UserClaims `json:"claims"`
	jwt.RegisteredClaims
}

type IAuthToken interface {
	SignToken() string
}

type IAdminToken interface {
	SignToken() string
}

type IApiKey interface {
	SignToken() string
}

func jwtTimeDurationCal(t int) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Now().Add(time.Duration(int64(t) * int64(math.Pow10(9)))))
}

func jwtTimeRepeatAdapter(t int64) *jwt.NumericDate {
	return jwt.NewNumericDate(time.Unix(t, 0))
}

func (a *basicAuth) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	tokenString, _ := token.SignedString(a.cfg.SecretKey())

	return tokenString
}

func (a *basicAdmin) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	tokenString, _ := token.SignedString(a.cfg.AdminKey())

	return tokenString
}

func (a *basicApiKey) SignToken() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, a.mapClaims)
	tokenString, _ := token.SignedString(a.cfg.ApiKey())

	return tokenString
}

func ParseToken(cfg config.IJwtConfig, tokenString string) (*basicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &basicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return cfg.SecretKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("invalid token")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		} else {
			return nil, fmt.Errorf("failed to parse token: %w", err)
		}
	}
	if claims, ok := token.Claims.(*basicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
}

func ParseAdminToken(cfg config.IJwtConfig, tokenString string) (*basicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &basicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return cfg.AdminKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("invalid token")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		} else {
			return nil, fmt.Errorf("failed to parse token: %w", err)
		}
	}
	if claims, ok := token.Claims.(*basicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
}

func ParseApiKey(cfg config.IJwtConfig, tokenString string) (*basicMapClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &basicMapClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return cfg.ApiKey(), nil
	})
	if err != nil {
		if errors.Is(err, jwt.ErrTokenMalformed) {
			return nil, fmt.Errorf("invalid token")
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, fmt.Errorf("token expired")
		} else {
			return nil, fmt.Errorf("failed to parse token: %w", err)
		}
	}
	if claims, ok := token.Claims.(*basicMapClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}
}

func RepeatToken(cfg config.IJwtConfig, claims *users.UserClaims, exp int64) string {
	obj := &basicAuth{
		cfg: cfg,
		mapClaims: &basicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "basicshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeRepeatAdapter(exp),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
	return obj.SignToken()
}

func NewAuth(tokenType TokenType, cfg config.IJwtConfig, claims *users.UserClaims) (IAuthToken, error) {
	switch tokenType {
	case Access:
		return newAccessToken(cfg, claims), nil
	case Refresh:
		return newRefreshToken(cfg, claims), nil
	case Admin:
		return newAdminToken(cfg), nil
	case ApiKey:
		return newApiKey(cfg), nil
	default:
		return nil, fmt.Errorf("invalid token type")
	}
}

func newAccessToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuthToken {
	return &basicAuth{
		cfg: cfg,
		mapClaims: &basicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "basicshop-api",
				Subject:   "access-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.AccessExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newRefreshToken(cfg config.IJwtConfig, claims *users.UserClaims) IAuthToken {
	return &basicAuth{
		cfg: cfg,
		mapClaims: &basicMapClaims{
			Claims: claims,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "basicshop-api",
				Subject:   "refresh-token",
				Audience:  []string{"customer", "admin"},
				ExpiresAt: jwtTimeDurationCal(cfg.RefreshExpiresAt()),
				NotBefore: jwt.NewNumericDate(time.Now()),
				IssuedAt:  jwt.NewNumericDate(time.Now()),
			},
		},
	}
}

func newAdminToken(cfg config.IJwtConfig) IAuthToken {
	return &basicAdmin{
		basicAuth: &basicAuth{
			cfg: cfg,
			mapClaims: &basicMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "basicshop-api",
					Subject:   "admin-token",
					Audience:  []string{"admin"},
					ExpiresAt: jwtTimeDurationCal(300),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}

func newApiKey(cfg config.IJwtConfig) IAuthToken {
	return &basicApiKey{
		basicAuth: &basicAuth{
			cfg: cfg,
			mapClaims: &basicMapClaims{
				Claims: nil,
				RegisteredClaims: jwt.RegisteredClaims{
					Issuer:    "basicshop-api",
					Subject:   "api-key",
					Audience:  []string{"admin", "customer"},
					ExpiresAt: jwt.NewNumericDate(time.Now().AddDate(2, 0, 0)),
					NotBefore: jwt.NewNumericDate(time.Now()),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
				},
			},
		},
	}
}
