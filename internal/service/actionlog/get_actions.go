package actionlog

import (
	"context"
	"log"
	"time"

	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (s *srvc) GetActions(ctx context.Context, req *model.ActionSearch) (<-chan *model.Action, error) {
	log.Println("GetActions request received:", req)

	out := make(chan *model.Action)

	go func() {
		defer close(out)

		actions := []*model.Action{
			{
				ID: 1,
				Info: &model.ActionInfo{
					UserID:     100,
					ActionType: model.ActionTypeBuy,
					Timestamp:  time.Now(),
					Details:    map[string]any{"item": "book"},
				},
			},
			{
				ID: 2,
				Info: &model.ActionInfo{
					UserID:     200,
					ActionType: model.ActionTypeRefund,
					Timestamp:  time.Now(),
					Details:    map[string]any{"item": "pen"},
				},
			},
		}

		for _, action := range actions {
			select {
			case <-ctx.Done():
				return
			case out <- action:
			}
		}
	}()

	return out, nil
}
