package service

// GenerateToken 生成token
import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// GenerateToken generates a new JWT token for a given user ID
func GenerateToken(userID int64) (string, error) {
	// Get the secret key from the environment
	SECRET_KEY, err := getSecretKey()
	if err != nil {
		return "", err
	}
	// Create the JWT claims, which includes the user ID and expiry time
	claims := &jwt.MapClaims{
		"userId": userID,
		"exp":    time.Now().Add(24 * time.Hour).Unix(), // Token expires after 24 hours
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(SECRET_KEY))
}

func getSecretKey() (string, error) {
	err := godotenv.Load()
	if err != nil {
		return "", err
	}
	return os.Getenv("SECRET_KEY"), nil
}
