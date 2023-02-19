package dblayer

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

var jwtKey = []byte("key_Grpc_66#6$1*_Sari@")

type Claims struct {
	userName string
	jwt.RegisteredClaims
}

type IGenerateTokenService interface {
	CreateToken(IssuerName string) string
	CreateHash256(Token string) string
}
type GenerateTokenServiceStruct struct {
}

func NewGenerateTokenServiceStruct() IGenerateTokenService {

	return &GenerateTokenServiceStruct{}
}

func (tsToken *GenerateTokenServiceStruct) CreateToken(IssuerName string) string {
	expirTime := time.Now().Add(24 * time.Hour)

	Uuid := uuid.New()
	Guid := strings.Replace(Uuid.String(), "-", "", -1)

	claims := &Claims{
		userName: IssuerName,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    IssuerName,
			Subject:   "Grpc",
			ID:        Guid,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expirTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "not create token"
	}

	return tokenString
}

func (tsToken *GenerateTokenServiceStruct) CreateHash256(Token string) string {
	hash := sha256.New()
	hash.Write([]byte(Token))
	hashToken := hex.EncodeToString(hash.Sum(nil)) // or ==>  hash := sha256.Sum256([]byte(Token))
	return hashToken
}
