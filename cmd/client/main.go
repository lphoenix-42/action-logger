package main

import (
	"context"
	"fmt"
	"log"

	"connectrpc.com/connect"
	"github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/gen/actionlog/v1/actionlog_v1connect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
)

func main() {
	ctx := context.Background()

	client := actionlog_v1connect.NewActionLogServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)

	// --- 1. LogAction
	log.Printf("Calling LogAction...\n")
	details, _ := structpb.NewStruct(map[string]interface{}{
		"item": "notebook",
	})
	logReq := &actionlog_v1.LogActionRequest{
		Info: &actionlog_v1.ActionInfo{
			UserId:     123,
			ActionType: actionlog_v1.ActionType_ACTION_TYPE_BUY,
			Timestamp:  timestamppb.Now(),
			Details:    details,
		},
	}
	logRes, err := client.LogAction(ctx, connect.NewRequest(logReq))
	if err != nil {
		log.Fatalf("LogAction error: %v", err)
	}
	fmt.Printf("LogAction response: ID = %d\n", logRes.Msg.Id)

	// --- 2. GetActions (stream)
	fmt.Println("\nStreaming GetActions...")
	streamReq := connect.NewRequest(&actionlog_v1.GetActionsRequest{})
	stream, err := client.GetActions(ctx, streamReq)
	if err != nil {
		log.Fatalf("GetActions error: %v", err)
	}
	for stream.Receive() {
		resp := stream.Msg()
		fmt.Printf("GetActions => Action ID: %d, User: %d\n", resp.Action.Id, resp.Action.Info.UserId)
	}
	if err := stream.Err(); err != nil {
		log.Fatalf("Stream receive error: %v", err)
	}

	// --- 3. WatchActions (stream)
	fmt.Println("\nStreaming WatchActions...")
	watchStream, err := client.WatchActions(ctx, connect.NewRequest(&actionlog_v1.WatchActionsRequest{}))
	if err != nil {
		log.Fatalf("WatchActions error: %v", err)
	}
	for watchStream.Receive() {
		resp := watchStream.Msg()
		fmt.Printf("Watch => Action ID: %d, Details: %v\n", resp.Action.Id, resp.Action.Info.Details.AsMap())
	}
	if err := watchStream.Err(); err != nil {
		log.Fatalf("Watch stream error: %v", err)
	}
}
