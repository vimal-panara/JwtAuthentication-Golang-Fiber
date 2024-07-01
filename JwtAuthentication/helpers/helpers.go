package helpers

import (
	"JwtAuthentication/handlers"
	"JwtAuthentication/models"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var SecreteKey = os.Getenv("JWT_SECRETE_KEY")
var Password_enc_key = os.Getenv("PASSWORD_ENC_KEY")
var Aes_iv = os.Getenv("AES_IV")

func GenerateJwtToken(user models.User) (string, string, error) {

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"mobile":    user.Mobile,
			"email":     user.Email,
			"exp":       time.Now().Add(time.Hour * 24).Unix(),
		},
	)
	// fmt.Println("Jwt Token secrete key ", secreteKey)
	tokenString, err := token.SignedString([]byte(SecreteKey))
	if err != nil {
		return "", "", err
	}

	refToken := jwt.NewWithClaims(
		jwt.SigningMethodHS256,
		jwt.MapClaims{
			"firstName": user.FirstName,
			"lastName":  user.LastName,
			"mobile":    user.Mobile,
			"email":     user.Email,
			"exp":       time.Now().Add(time.Hour * 24 * 3).Unix(),
		},
	)
	refTokenString, err := refToken.SignedString([]byte(SecreteKey))
	if err != nil {
		return "", "", err
	}
	return tokenString, refTokenString, nil
}

func LoadEncryptionKeys() {

	SecreteKey = os.Getenv("JWT_SECRETE_KEY")
	Password_enc_key = os.Getenv("PASSWORD_ENC_KEY")
	Aes_iv = os.Getenv("AES_IV")

	if SecreteKey == "" {
		SecreteKey = "dfkhlajfhlkasdfhjsenfhejrnskfnkjsgkjsdnfksdnfkjsdhfjdhfkjdsncs"
	}
	if Password_enc_key == "" {
		Password_enc_key = "4663fc6d9a154ec580a7f75d89555fdf"
	}
	if Aes_iv == "" {
		Aes_iv = "3b734604e2574a3cac0813bd257ce84f"
	}
}

func GetEncryptedPassword(password string) (string, error) {

	// fmt.Println("Password Encryption Key ", password_enc_key)
	block, err := aes.NewCipher([]byte(Password_enc_key))
	if err != nil {
		return "", err
	}
	// fmt.Println("New Cipher Block ", block)
	plainText := []byte(password)
	cfb := cipher.NewCFBEncrypter(block, []byte(Aes_iv)[:aes.BlockSize])
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func GetDecryptedPassword(encryptedPassword string) (string, error) {

	// fmt.Println("Password Decryption Key ", password_enc_key)
	block, err := aes.NewCipher([]byte(Password_enc_key))
	if err != nil {
		return "", err
	}
	cipherText, err := base64.StdEncoding.DecodeString(encryptedPassword)
	if err != nil {
		return "", err
	}
	cfb := cipher.NewCFBDecrypter(block, []byte(Aes_iv)[:aes.BlockSize])
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)
	return string(plainText), nil
}

func ValidateJwtToken(token string) error {
	jwtToken, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecreteKey), nil
	})
	if err != nil {
		return err
	}

	if !jwtToken.Valid {
		return errors.New("invalid token recieved")
	}

	claims := jwtToken.Claims.(jwt.MapClaims)
	exp := claims["exp"].(float64)

	if exp < float64(time.Now().Unix()) {
		return errors.New("invalid token")
	}

	return nil
}

func UpdateUserTokens(token string) error {

	jwtToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecreteKey), nil
	})

	claims := jwtToken.Claims.(jwt.MapClaims)

	user := models.User{
		FirstName: claims["firstName"].(string),
		LastName:  claims["lastName"].(string),
		Email:     claims["email"].(string),
		Mobile:    claims["mobile"].(string),
	}

	newToken, newRefToken, _ := GenerateJwtToken(user)

	user.Token = newToken
	user.RefreshToken = newRefToken
	user.UpdatedAt = primitive.NewDateTimeFromTime(time.Now())

	handlers.UpdateUserTokens(user, primitive.NewObjectID())

	return nil
}

func GetEmailMobileFromToken(token string) (models.User, error) {
	jwtToken, _ := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(SecreteKey), nil
	})

	claims := jwtToken.Claims.(jwt.MapClaims)

	user := models.User{
		FirstName: claims["firstName"].(string),
		LastName:  claims["lastName"].(string),
		Email:     claims["email"].(string),
		Mobile:    claims["mobile"].(string),
	}

	return user, nil
}
