package user

import (
	"github.com/f1xend/auth/internal/repository"
	"github.com/f1xend/auth/internal/service"
	"github.com/f1xend/platform-common/pkg/db"
)

var _ service.UserService = (*serv)(nil)

type serv struct {
	userRepository repository.UserRepository
	txManager      db.TxManager
}

func NewService(userRepository repository.UserRepository, txManager db.TxManager) service.UserService {
	return &serv{
		userRepository: userRepository,
		txManager:      txManager,
	}
}
