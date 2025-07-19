package actionlog

import (
	"context"
	"log"
	"time"

	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (s *srvc) WatchActions(ctx context.Context) (<-chan *model.Action, error) {
	log.Println("WatchActions started")

	ch := make(chan *model.Action)

	go func() {
		defer close(ch)

		for i := int64(1); i <= 3; i++ {
			ch <- &model.Action{
				ID: i,
				Info: &model.ActionInfo{
					UserID:     500 + i,
					ActionType: model.ActionTypeBuy,
					Timestamp:  time.Now(),
					Details:    map[string]any{"demo": i},
				},
			}
			time.Sleep(500 * time.Millisecond)
		}

		log.Println("WatchActions ended")
	}()

	return ch, nil
}
