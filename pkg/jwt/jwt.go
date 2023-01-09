package jwtToken

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v4"
)

// fungsi untuk membuat token
func GenerateToken(claims *jwt.MapClaims) (string, error) {
	// membuat token baru
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// men-sign token yang sudah dibuat + secret key
	webtoken, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", err
	}

	return webtoken, nil
}

// fungsi untuk verifikasi token
func VerifyToken(tokenString string) (*jwt.Token, error) {
	// menghilangkan string `Bearer ` pada authorization
	tokenString = strings.Replace(tokenString, "Bearer ", "", -1)

	// memparsing token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// jika method tidak valid, maka kirim error sebagai kembalian dari fungsi callback
		if _, isValid := token.Method.(*jwt.SigningMethodHMAC); !isValid {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		// jika method valid, maka kirim secret key sebagai kembalian dari fungsi callback
		return []byte(os.Getenv("SECRET_KEY")), nil
	})

	// jika parsing gagal, kembalikan error
	if err != nil {
		return nil, err
	}

	// jika parsing berhasil, kemabalikan token
	return token, nil
}

// fungsi untuk decoding token
func DecodeToken(tokenString string) (jwt.MapClaims, error) {
	// memverifikasi token menggunakan fungsi yang dibuat diatas
	token, err := VerifyToken(tokenString)
	if err != nil {
		return nil, err
	}

	// mengambil informasi payload/claims pada jwt
	claims, isOK := token.Claims.(jwt.MapClaims)
	// pengecekan valid-tidaknya claims tersebut, jika tidak valid maka kembalikan error
	if !isOK || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	// jika valid, kembalikan nilai payload/claims
	return claims, nil
}
