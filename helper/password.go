package helper

import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) string {
	pass := []byte(password)
	hashedPassword, _ := bcrypt.GenerateFromPassword(pass, bcrypt.DefaultCost)
	return string(hashedPassword)
}
