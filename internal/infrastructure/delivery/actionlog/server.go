package actionlog

type Server struct{}

func New() *Server {
	return &Server{}
}

// func New(actionlogService service.ActionlogService) *Server {
// 	return &Server{
// 		service: actionlogService,
// 	}
// }
