package utils

import (
	"github.com/Alf-Grindel/clide/pkg/constants"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

func addSalt(password string) string {
	var builder strings.Builder
	builder.WriteString(password)
	builder.WriteString(constants.Salt)
	return builder.String()
}

func GeneratePassword(password string) (string, error) {
	saltPassword := addSalt(password)

	b, err := bcrypt.GenerateFromPassword([]byte(saltPassword), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ComparePassword(hashPassword, password string) bool {
	saltPassword := addSalt(password)

	err := bcrypt.CompareHashAndPassword([]byte(hashPassword), []byte(saltPassword))
	return err == nil
}
