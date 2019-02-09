package jwt

import (
	"github.com/alamin-mahamud/container-management-system/user-service/pkg/repository"
	pb "github.com/alamin-mahamud/container-management-system/user-service/proto/user"
	"github.com/dgrijalva/jwt-go"
)

var (
	// Define a secure key string used
	// as a salt when hashing our tokens.
	// Please make your own way more secure than this,
	// use a randomly generated md5 hash or something.
	key = []byte("mySuperSecretKeyLol")
)

// CustomClaims is our custom metadata, which will be hashed
// and sent as the second segment in our JWT
type CustomClaims struct {
	User *pb.User
	jwt.StandardClaims
}

// Authable - ...
type Authable interface {
	Encode(user *pb.User) (string, error)
	Decode(token string) (*CustomClaims, error)
}

type TokenService struct {
	Repo repository.Repository
}

// Encode Token
func (srv *TokenService) Encode(user *pb.User) (string, error) {
	// Create the claims
	claims := CustomClaims{
		user,
		jwt.StandardClaims{
			ExpiresAt: 86400,
			Issuer:    "go.micro.srv.user",
		},
	}

	// Create Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token and return
	return token.SignedString(key)
}

// Decode a token string into a token object
func (srv *TokenService) Decode(token string) (*CustomClaims, error) {
	// Parse the token
	tokenType, err := jwt.ParseWithClaims(token, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})

	// validate the token and return the custom claims
	if claims, ok := tokenType.Claims.(*CustomClaims); ok && tokenType.Valid {
		return claims, nil
	}

	return nil, err

}
