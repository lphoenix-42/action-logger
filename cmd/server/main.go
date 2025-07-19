package main

import (
	"context"
	"log"
	"net/http"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/gen/actionlog/v1/actionlog_v1connect"
)

type Server struct{}

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

func main() {
	al := &Server{}
	mux := http.NewServeMux()
	path, handler := actionlog_v1connect.NewActionLogServiceHandler(al)
	mux.Handle(path, handler)
	log.Println("Up and running")
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
