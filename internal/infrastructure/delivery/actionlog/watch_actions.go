package actionlog

import (
	"context"
	"log"

	"connectrpc.com/connect"
	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) WatchActions(
	ctx context.Context,
	req *connect.Request[actionlogv1.WatchActionsRequest],
	stream *connect.ServerStream[actionlogv1.WatchActionsResponse],
) error {
	log.Println("WatchActions started")

	for i := int64(1); i <= 3; i++ {
		act := &actionlogv1.Action{
			Id: i,
			Info: &actionlogv1.ActionInfo{
				UserId:     500 + i,
				ActionType: actionlogv1.ActionType_ACTION_TYPE_BUY,
				Timestamp:  timestamppb.Now(),
				Details:    &structpb.Struct{Fields: map[string]*structpb.Value{"demo": structpb.NewNumberValue(float64(i))}},
			},
		}
		if err := stream.Send(&actionlogv1.WatchActionsResponse{Action: act}); err != nil {
			return err
		}
	}

	log.Println("WatchActions ended")
	return nil
}
