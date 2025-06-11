package utils

import (
	"fmt"
	"regexp"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey = []byte("HJGJH!a`#@!-@-`~@12-901asdAZw")

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func CheckPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}

func CreateToken(userClaims []map[string]interface{}) (string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"userClaims": userClaims,
			"exp":        time.Now().Add(time.Hour * 24).Unix(),
		})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func isValidEmail(email string) bool {

	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}
