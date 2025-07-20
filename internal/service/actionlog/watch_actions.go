package actionlog

import (
	"context"
	"fmt"

	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (s *srvc) WatchActions(ctx context.Context) (<-chan *model.Action, error) {
	notifyChan, err := s.repo.ListenNotifications(ctx, "new_action")
	if err != nil {
		return nil, err
	}

	ch := make(chan *model.Action)
	go func() {
		defer close(ch)

		for payload := range notifyChan {
			var id int64
			if _, err := fmt.Sscanf(payload, "%d", &id); err != nil {
				continue
			}

			action, err := s.repo.GetByID(ctx, id)
			if err != nil || action == nil {
				continue
			}

			ch <- action
		}
	}()

	return ch, nil
}
