package apperrors

import (
	"errors"
)

var ErrEmailAlreadyExists = errors.New("email já cadastrado")
var ErrInvalidEmailOrPassword = errors.New("email ou senha inválidos")
var ErrUserNotFound = errors.New("usuário não encontrado")
var ErrNotFound = errors.New("recurso não encontrado")
var ErrInvalidRefreshToken = errors.New("refresh token inválido ou expirado")
