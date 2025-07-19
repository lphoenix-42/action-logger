package actionlog

import (
	"context"

	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (s *srvc) LogAction(ctx context.Context, info *model.ActionInfo) (int64, error) {
	return 454, nil
}
