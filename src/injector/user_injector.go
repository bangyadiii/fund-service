//go:build wireinject
// +build wireinject

package injector

import (
	"backend-crowdfunding/src/repository"
	"backend-crowdfunding/src/service"

	"github.com/google/wire"
)

func InitializeUserService() service.UserService {
	wire.Build(repository.NewUserRepository, service.NewUserService)
	return nil
}
