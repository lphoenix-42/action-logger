package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"net/http"

	"connectrpc.com/connect"
	"github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/gen/actionlog/v1/actionlog_v1connect"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
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

	// --- 2. GetActions (stream, no filter)
	counter := 0
	fmt.Println("\nStreaming GetActions (no filter)...")
	streamReq := connect.NewRequest(&actionlog_v1.GetActionsRequest{
		DetailsFilters: []*actionlog_v1.JsonFilter{},
	})
	getStream, err := client.GetActions(ctx, streamReq)
	if err != nil {
		log.Fatalf("GetActions error: %v", err)
	}
	for getStream.Receive() {
		resp := getStream.Msg()
		PrintAction(resp.Action)
		counter++
	}
	if err := getStream.Err(); err != nil {
		log.Fatalf("Stream receive error: %v", err)
	}
	fmt.Println("\nTotal: ", counter)

	// --- 3. GetActions (stream, with filter)
	counter = 0
	filterId := int64(42)
	fmt.Println("\nStreaming GetActions (with filter)...")
	streamReq = connect.NewRequest(&actionlog_v1.GetActionsRequest{
		UserId: &filterId,
		DetailsFilters: []*actionlog_v1.JsonFilter{
			{
				Path:  []string{"catalog", "item"},
				Value: "Notebook",
			},
			{
				Path:  []string{"discount"},
				Value: "true",
			},
		},
	})
	getFilterStream, err := client.GetActions(ctx, streamReq)
	if err != nil {
		log.Fatalf("GetActions error: %v", err)
	}
	for getFilterStream.Receive() {
		resp := getFilterStream.Msg()
		PrintAction(resp.Action)
		counter++
	}
	if err := getFilterStream.Err(); err != nil {
		log.Fatalf("Stream receive error: %v", err)
	}
	fmt.Println("\nTotal: ", counter)

	// --- 4. WatchActions (stream)
	fmt.Println("\nStreaming WatchActions...")
	watchStream, err := client.WatchActions(ctx, connect.NewRequest(&actionlog_v1.WatchActionsRequest{}))
	if err != nil {
		log.Fatalf("WatchActions error: %v", err)
	}
	for watchStream.Receive() {
		resp := watchStream.Msg()
		PrintAction(resp.Action)
	}
	if err := watchStream.Err(); err != nil {
		log.Fatalf("Watch stream error: %v", err)
	}
}

func PrintAction(act *actionlog_v1.Action) {
	if act == nil || act.Info == nil {
		fmt.Println("Empty action or missing info")
		return
	}

	timestamp := "nil"
	if act.Info.Timestamp != nil {
		timestamp = act.Info.Timestamp.AsTime().Format(time.RFC3339)
	}

	details := "{}"
	if act.Info.Details != nil {
		if b, err := json.MarshalIndent(act.Info.Details.AsMap(), "", "  "); err == nil {
			details = string(b)
		}
	}

	fmt.Printf(`Action id: %d, user: %d, type: %s, timestamp: %s
Details: %s

`, act.Id, act.Info.UserId, act.Info.ActionType.String(), timestamp, details)
}
