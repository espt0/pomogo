package apperrors

import (
	"errors"
)

var ErrEmailAlreadyExists = errors.New("email já cadastrado")
var ErrInvalidEmailOrPassword = errors.New("email ou senha inválidos")

var ErrRefreshInvalido = errors.New("refresh token inválido ou expirado")
