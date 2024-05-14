package cryptor

import "golang.org/x/crypto/bcrypt"

func HashPassword(txt string) (string, error) {

	bytes := make([]byte, 100)
	defer func() {
		bytes = nil
	}()

	copy(bytes, txt)

	result, err := bcrypt.GenerateFromPassword(bytes, 12)
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func CompareSecureToPlain(secure string, plain string) bool {

	secBytes := make([]byte, 100)
	plainBytes := make([]byte, 100)

	defer func() {
		secBytes = nil
		plainBytes = nil
	}()

	copy(secBytes, secure)
	copy(plainBytes, plain)

	return bcrypt.CompareHashAndPassword(secBytes, plainBytes) == nil
}