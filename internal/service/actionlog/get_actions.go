package actionlog

import (
	"context"

	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (s *srvc) GetActions(ctx context.Context, req *model.ActionSearch) (<-chan *model.Action, error) {
	actions, err := s.repo.QueryActions(ctx, req)
	if err != nil {
		return nil, err
	}

	ch := make(chan *model.Action, len(actions))
	go func() {
		defer close(ch)
		count := 0
		for _, action := range actions {
			if count >= 100 {
				break
			}
			ch <- action
			count++
		}
	}()

	return ch, nil
}
