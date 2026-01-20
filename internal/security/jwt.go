package security

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/razedwell/go-hand/internal/model"
	"github.com/razedwell/go-hand/internal/platform/cache"
	"github.com/razedwell/go-hand/internal/repository/token"
	"github.com/razedwell/go-hand/internal/transport/http/helpers"
)

type JWTClaims struct {
	UserID int64  `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

type JWTManager struct {
	accessSecret  []byte
	refreshSecret []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
	repo          token.Repository   // Your Postgres Repo
	redis         *cache.RedisClient // For Logout Blacklist
}

func NewJWTManager(accessSecret, refreshSecret string, accessExpiry, refreshExpiry time.Duration, repo token.Repository, rdb *cache.RedisClient) *JWTManager {
	return &JWTManager{
		accessSecret:  []byte(accessSecret),
		refreshSecret: []byte(refreshSecret),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		repo:          repo,
		redis:         rdb,
	}
}

func (j *JWTManager) hashToken(token string) string {
	h := sha256.New()
	h.Write([]byte(token))
	return hex.EncodeToString(h.Sum(nil))
}

func (j *JWTManager) IsBlacklisted(ctx context.Context, accessTokenStr string) bool {
	n, _ := j.redis.Client.Exists(ctx, accessTokenStr).Result()
	return n > 0
}

func (j *JWTManager) Verify(tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.accessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}

func (j *JWTManager) GenerateTokenPair(userID int64, role string) (string, string, error) {
	accessClaims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(helpers.GetCurrentTimeStampUTC().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(helpers.GetCurrentTimeStampUTC()),
		},
	}
	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(j.accessSecret)
	if err != nil {
		return "", "", err
	}

	// 2. Generate Refresh Token
	refreshExpiryTime := helpers.GetCurrentTimeStampUTC().Add(j.refreshExpiry)
	refreshClaims := &jwt.RegisteredClaims{
		Subject:   strconv.FormatInt(userID, 10),
		ExpiresAt: jwt.NewNumericDate(refreshExpiryTime),
		IssuedAt:  jwt.NewNumericDate(helpers.GetCurrentTimeStampUTC()),
	}
	refreshToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString(j.refreshSecret)
	if err != nil {
		return "", "", err
	}

	hash := j.hashToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	err = j.repo.CreateRefreshToken(context.Background(), &model.RefreshToken{
		UserID:    userID,
		TokenHash: hash,
		ExpiresAt: refreshExpiryTime,
	})

	return accessToken, refreshToken, err
}

func (j *JWTManager) generateAccessToken(userID int64, role string) (string, error) {
	accessClaims := &JWTClaims{
		UserID: userID,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(helpers.GetCurrentTimeStampUTC().Add(j.accessExpiry)),
			IssuedAt:  jwt.NewNumericDate(helpers.GetCurrentTimeStampUTC()),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString(j.accessSecret)
}

func (j *JWTManager) RefreshAccessToken(ctx context.Context, refreshTokenStr string) (string, error) {
	// 1. Verify Refresh Token Signature
	token, err := jwt.Parse(refreshTokenStr, func(t *jwt.Token) (interface{}, error) {
		return j.refreshSecret, nil
	})
	if err != nil || !token.Valid {
		return "", errors.New("invalid refresh token")
	}

	// 2. Check DB for the hash
	hash := j.hashToken(refreshTokenStr)
	storedToken, err := j.repo.GetRefreshToken(ctx, hash)
	if err != nil {
		return "", errors.New("refresh token not found")
	}

	// 3. Security Checks
	if storedToken.RevokedAt != nil {
		return "", errors.New("refresh token was revoked")
	}
	if helpers.GetCurrentTimeStampUTC().After(storedToken.ExpiresAt) {
		return "", errors.New("refresh token expired")
	}

	// 4. Generate NEW Access Token (Keep the user logged in)
	// You might want to fetch the latest Role from UserRepo here
	return j.generateAccessToken(storedToken.UserID, "user")
}

func (j *JWTManager) BlacklistTokens(ctx context.Context, accessTokenStr string, refreshTokenStr string) error {
	// Blacklist Access Token in Redis
	accessToken, err := jwt.ParseWithClaims(accessTokenStr, &JWTClaims{}, func(t *jwt.Token) (interface{}, error) {
		return j.accessSecret, nil
	})
	if err == nil && accessToken.Valid {
		claims := accessToken.Claims.(*JWTClaims)
		expiry := claims.ExpiresAt.Time.Sub(helpers.GetCurrentTimeStampUTC())
		j.redis.Client.Set(ctx, accessTokenStr, "blacklisted", expiry)
	}

	// Revoke Refresh Token in DB
	hash := j.hashToken(refreshTokenStr)
	return j.repo.RevokeRefreshToken(ctx, hash)
}
