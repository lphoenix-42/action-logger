package service

import (
	"context"

	"github.com/lphoenix-42/action-logger/internal/service/model"
)

type ActionlogService interface {
	GetActions(ctx context.Context, req *model.ActionSearch) (<-chan *model.Action, error)
	LogAction(ctx context.Context, info *model.ActionInfo) (int64, error)
	WatchActions(ctx context.Context) (<-chan *model.Action, error)
}
