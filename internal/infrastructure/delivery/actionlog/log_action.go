package actionlog

import (
	"context"

	"connectrpc.com/connect"
	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/converter"
)

func (s *Server) LogAction(
	ctx context.Context,
	req *connect.Request[actionlogv1.LogActionRequest],
) (*connect.Response[actionlogv1.LogActionResponse], error) {
	id, err := s.service.LogAction(ctx, converter.ActionInfoFromDescToModel(req.Msg.Info))
	res := connect.NewResponse(&actionlogv1.LogActionResponse{
		Id: id,
	})
	return res, err
}
