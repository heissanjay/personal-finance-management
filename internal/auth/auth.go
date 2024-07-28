package auth

import (
	"context"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/heissanjay/personal-finance-management/internal/config"
	"github.com/heissanjay/personal-finance-management/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

var (
	jwtKey = []byte(config.Config.JWTSecret)
)

type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

type AuthService struct {
	Storage storage.DataStore
}

func NewAuthService(storage storage.DataStore) *AuthService {
	return &AuthService{Storage: storage}
}

func (a *AuthService) RegisterUser(username, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	newUser := storage.User{
		Username: username,
		Password: string(hashedPassword),
	}

	return a.Storage.SaveUser(context.Background(), newUser)
}

func (a *AuthService) LoginUser(username, password string) (string, error) {
	var registeredUser storage.User
	registeredUser, err := a.Storage.GetUserByUsername(context.Background(), username)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(registeredUser.Password), []byte(password)); err != nil {
		return "", err
	}

	expirationTime := time.Now().Add(24 * time.Hour)

	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
