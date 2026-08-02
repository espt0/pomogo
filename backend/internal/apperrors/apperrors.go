package apperrors

import (
	"errors"
)

var ErrEmailAlreadyExists = errors.New("email já cadastrado")
var ErrInvalidEmailOrPassword = errors.New("email ou senha inválidos")
var ErrUserNotFound = errors.New("usuário não encontrado")
var ErrNotFound = errors.New("recurso não encontrado")
var ErrInvalidRefreshToken = errors.New("refresh token inválido ou expirado")
var ErrArchivedTask = errors.New("não é possível iniciar sessão em tarefa arquivada")
var ErrTaskCompleted = errors.New("tarefa já concluída. Reabra-a antes de iniciar uma sessão.")
var ErrSessionAlreadyActive = errors.New("sessão já ativa")
