package actionlog

import (
	"github.com/lphoenix-42/action-logger/internal/infrastructure/repository"
	"github.com/lphoenix-42/action-logger/internal/service"
)

type srvc struct {
	repo repository.ActionlogRepository
}

func New(actionlogRepository repository.ActionlogRepository) service.ActionlogService {
	return &srvc{
		repo: actionlogRepository,
	}
}
