package notifier

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/repository"
)

type notifier struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) repository.Notifier {
	return &notifier{db: db}
}

func (n *notifier) Listen(ctx context.Context, channel string) (<-chan string, error) {
	conn, err := n.db.Acquire(ctx)
	if err != nil {
		return nil, err
	}

	if _, err := conn.Exec(ctx, `LISTEN "`+channel+`"`); err != nil {
		conn.Release()
		return nil, err
	}

	notifyChan := make(chan string)
	go func() {
		defer conn.Release()
		defer close(notifyChan)

		for {
			notification, err := conn.Conn().WaitForNotification(ctx)
			if err != nil {
				return
			}
			notifyChan <- notification.Payload
		}
	}()

	return notifyChan, nil
}

func (n *notifier) Notify(ctx context.Context, channel, payload string) error {
	_, err := n.db.Exec(ctx, fmt.Sprintf(`NOTIFY "%s", '%s'`, channel, payload))
	return err
}
