package authenticator

import (
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/ed25519"
)

type JwtAuthenticator struct {
	privateKey      ed25519.PrivateKey
	publicKey       ed25519.PublicKey
	issuer          string
	accessTokenTTL  time.Duration
	refreshTokenTTL time.Duration
}

type TokenPair struct {
	AccessToken           string
	RefreshToken          string
	RefreshTokenExpiresAt time.Time
}

func New(privateKey, publicKey []byte, accessTokenTTL, refreshTokenTTL time.Duration, issuer string) (*JwtAuthenticator, error) {
	priv, err := loadEd25519PrivateKey(privateKey)
	if err != nil {
		return nil, err
	}

	pub, err := loadEd25519PublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	return &JwtAuthenticator{
		privateKey:      priv,
		publicKey:       pub,
		issuer:          issuer,
		accessTokenTTL:  accessTokenTTL,
		refreshTokenTTL: refreshTokenTTL,
	}, nil
}

func loadEd25519PrivateKey(privateKey []byte) (ed25519.PrivateKey, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}

	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	privateKey, ok := key.(ed25519.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("invalid private key type")
	}

	return privateKey, nil
}

func loadEd25519PublicKey(publicKey []byte) (ed25519.PublicKey, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, fmt.Errorf("invalid PEM data")
	}

	key, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	publicKey, ok := key.(ed25519.PublicKey)
	if !ok {
		return nil, fmt.Errorf("invalid public key type")
	}

	return publicKey, nil
}

func (j *JwtAuthenticator) IssueTokens(userID uuid.UUID) (*TokenPair, error) {
	now := time.Now()

	accessToken, err := j.generateAccessToken(userID, now)
	if err != nil {
		return nil, err
	}

	refreshToken, err := j.generateRefreshToken()
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:           accessToken,
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: now.Add(j.refreshTokenTTL),
	}, nil
}

func (j *JwtAuthenticator) ValidateAccessToken(tokenStr string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
		return j.publicKey, nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrTokenInvalidClaims
}

func (j *JwtAuthenticator) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr := c.GetHeader("Authorization")
		if tokenStr == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			return
		}

		parts := strings.SplitN(tokenStr, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid Authorization header format",
			})
			return
		}

		token := parts[1]
		if _, err := j.ValidateAccessToken(token); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token",
			})
			return
		}

		c.Next()
	}
}

func (j *JwtAuthenticator) generateAccessToken(userID uuid.UUID, now time.Time) (string, error) {
	claims := jwt.RegisteredClaims{
		Issuer:    j.issuer,
		Subject:   userID.String(),
		ExpiresAt: jwt.NewNumericDate(now.Add(j.accessTokenTTL)),
		IssuedAt:  jwt.NewNumericDate(now),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodEdDSA, claims)
	return token.SignedString(j.privateKey)
}

func (j *JwtAuthenticator) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(bytes), nil
}
