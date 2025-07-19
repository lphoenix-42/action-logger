package actionlog

import (
	"github.com/lphoenix-42/action-logger/internal/service"
)

type srvc struct{}

func New() service.ActionlogService {
	return &srvc{}
}
