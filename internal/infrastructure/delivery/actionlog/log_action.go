package actionlog

import (
	"context"

	"connectrpc.com/connect"
	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
)

func (s *Server) LogAction(
	ctx context.Context,
	req *connect.Request[actionlogv1.LogActionRequest],
) (*connect.Response[actionlogv1.LogActionResponse], error) {
	res := connect.NewResponse(&actionlogv1.LogActionResponse{
		Id: 545,
	})
	res.Header().Set("ActionLogger-Version", "v1")
	return res, nil
}
