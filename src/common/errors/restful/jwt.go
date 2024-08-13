package restful

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

func HanldeJwtError(err error) error {
	if errors.Is(err, jwt.ErrInvalidKey) || errors.Is(err, jwt.ErrSignatureInvalid) {
		return errors.New("token is invalid")
	}

	if errors.Is(err, jwt.ErrTokenMalformed) {
		return errors.New("token is malformed")
	}

	if errors.Is(err, jwt.ErrTokenExpired) {
		return errors.New("token is expired")
	}

	return nil
}
