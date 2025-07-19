package actionlog

import (
	"context"
	"log"

	"connectrpc.com/connect"
	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Server) GetActions(
	ctx context.Context,
	req *connect.Request[actionlogv1.GetActionsRequest],
	stream *connect.ServerStream[actionlogv1.GetActionsResponse],
) error {
	log.Println("GetActions request received:", req.Msg)

	dummy := []actionlogv1.Action{
		{
			Id: 1,
			Info: &actionlogv1.ActionInfo{
				UserId:     100,
				ActionType: actionlogv1.ActionType_ACTION_TYPE_BUY,
				Timestamp:  timestamppb.Now(),
				Details:    &structpb.Struct{Fields: map[string]*structpb.Value{"item": structpb.NewStringValue("book")}},
			},
		},
		{
			Id: 2,
			Info: &actionlogv1.ActionInfo{
				UserId:     200,
				ActionType: actionlogv1.ActionType_ACTION_TYPE_REFUND,
				Timestamp:  timestamppb.Now(),
				Details:    &structpb.Struct{Fields: map[string]*structpb.Value{"item": structpb.NewStringValue("pen")}},
			},
		},
	}

	for _, act := range dummy {
		if err := stream.Send(&actionlogv1.GetActionsResponse{Action: &act}); err != nil {
			return err
		}
	}
	return nil
}
