package actionlog

import (
	"context"

	"connectrpc.com/connect"
	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/converter"
)

func (s *Server) GetActions(
	ctx context.Context,
	req *connect.Request[actionlogv1.GetActionsRequest],
	stream *connect.ServerStream[actionlogv1.GetActionsResponse],
) error {
	actionSearch := converter.FromGetActionsRequestToActionsSearchModel(req.Msg)

	actions, err := s.service.GetActions(ctx, actionSearch)
	if err != nil {
		return connect.NewError(connect.CodeInternal, err)
	}

	for act := range actions {
		desc := converter.ActionFromModelToDesc(act)
		if err := stream.Send(&actionlogv1.GetActionsResponse{Action: desc}); err != nil {
			return err
		}
	}
	return nil
}
