package repository

import (
	"context"
	"github.com/lphoenix-42/action-logger/internal/service/model"
)

type ActionlogRepository interface {
	GetByID(ctx context.Context, id int64) (*model.Action, error)
	InsertAction(ctx context.Context, info *model.ActionInfo) (int64, error)
	ListenNotifications(ctx context.Context, channel string) (<-chan string, error)
	QueryActions(ctx context.Context, req *model.ActionSearch) ([]*model.Action, error)
}
