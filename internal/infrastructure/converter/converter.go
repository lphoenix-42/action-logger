package converter

import (
	"time"

	actionlogv1 "github.com/lphoenix-42/action-logger/gen/actionlog/v1"
	"github.com/lphoenix-42/action-logger/internal/service/model"

	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ActionInfoFromDescToModel(info *actionlogv1.ActionInfo) *model.ActionInfo {
	return &model.ActionInfo{
		UserID:     info.UserId,
		ActionType: model.ActionType(info.ActionType),
		Timestamp:  info.Timestamp.AsTime(),
		Details:    info.Details.AsMap(),
	}
}

func ActionFromModelToDesc(a *model.Action) *actionlogv1.Action {
	if a == nil || a.Info == nil {
		return nil
	}

	detailsStruct, _ := structpb.NewStruct(a.Info.Details)

	return &actionlogv1.Action{
		Id: a.ID,
		Info: &actionlogv1.ActionInfo{
			UserId:     a.Info.UserID,
			ActionType: actionlogv1.ActionType(a.Info.ActionType),
			Timestamp:  timestamppb.New(a.Info.Timestamp),
			Details:    detailsStruct,
		},
	}
}

func FromGetActionsRequestToActionsSearchModel(req *actionlogv1.GetActionsRequest) *model.ActionSearch {
	var fromTime, toTime *time.Time

	if req.FromTimestamp != nil {
		t := req.FromTimestamp.AsTime()
		fromTime = &t
	}

	if req.ToTimestamp != nil {
		t := req.ToTimestamp.AsTime()
		toTime = &t
	}

	var filters []model.JsonFilter
	for _, f := range req.DetailsFilters {
		if f == nil {
			continue
		}
		filters = append(filters, model.JsonFilter{
			Path:  f.Path,
			Value: f.Value,
		})
	}

	return &model.ActionSearch{
		UserID:         req.UserId,
		ActionTypes:    toModelActionTypes(req.ActionType),
		From:           fromTime,
		To:             toTime,
		DetailsFilters: filters,
	}
}

func toModelActionTypes(types []actionlogv1.ActionType) []model.ActionType {
	result := make([]model.ActionType, 0, len(types))
	for _, t := range types {
		result = append(result, model.ActionType(t))
	}
	return result
}
