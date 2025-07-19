package actionlog

import "github.com/lphoenix-42/action-logger/internal/service"

type Server struct {
	service service.ActionlogService
}

func New(actionlogService service.ActionlogService) *Server {
	return &Server{
		service: actionlogService,
	}
}
