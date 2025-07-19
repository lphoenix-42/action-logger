package actionlog

import (
	"context"

	"connectrpc.com/connect"
	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/converter"
)

func (s *Server) WatchActions(
	ctx context.Context,
	req *connect.Request[actionlogv1.WatchActionsRequest],
	stream *connect.ServerStream[actionlogv1.WatchActionsResponse],
) error {
	ch, err := s.service.WatchActions(ctx)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	for act := range ch {
		desc := converter.ActionFromModelToDesc(act)
		if err := stream.Send(&actionlogv1.WatchActionsResponse{Action: desc}); err != nil {
			return err
		}
	}
	return nil
}
