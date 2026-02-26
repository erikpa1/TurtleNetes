package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"
	"turtle/users"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var (
	JWT_SECRET      = []byte("CHANGE_ME_SUPER_SECRET") // env var in prod
	JWT_ISSUER      = "files-receiver"
	JWT_EXPIRATION  = 12 * time.Hour
	JWT_COOKIE_NAME = "docminer_token"
)

type Claims struct {
	UserID string `json:"uid"`
	Email  string `json:"email"`
	Role   string `json:"role"`

	jwt.RegisteredClaims
}

func CreateToken(ip, uid, email, role string) (string, error) {
	claims := Claims{
		UserID: uid,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    JWT_ISSUER,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(JWT_EXPIRATION)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	secret := GetJwtString(ip, email, uid)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}

func ParseToken(tokenStr, ip string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(
		tokenStr,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, errors.New("unexpected signing method")
			}

			return GetJwtString(
				ip,
				claims.Email,
				claims.UserID,
			), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

func TryLogin(ctx context.Context, email, password string) (*users.User, error) {
	user, err := users.UserWithPasswordExists(ctx, email, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func SetAuthCookie(c *gin.Context, token string) {
	c.SetCookie(
		JWT_COOKIE_NAME,
		token,
		int(JWT_EXPIRATION.Seconds()),
		"/",
		"",
		false, // true if HTTPS
		true,  // httpOnly
	)
}

func LoginHandler(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user, err := TryLogin(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ip := c.ClientIP()

	token, err := CreateToken(
		ip,
		user.Uid.Hex(),
		user.Email,
		user.Role,
	)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	SetAuthCookie(c, token)

	c.JSON(http.StatusOK, gin.H{
		"uid":  user.Uid.Hex(),
		"role": user.Role,
	})
}

func GetJwtString(ip, email, uid string) []byte {
	date := time.Now().Format("2006-01-02") // YYYY-MM-DD

	base := fmt.Sprintf(
		"%s|%s_%s^%s",
		ip,
		email,
		uid,
		date,
	)

	hash := sha256.Sum256([]byte(base))
	return []byte(hex.EncodeToString(hash[:]))
}
