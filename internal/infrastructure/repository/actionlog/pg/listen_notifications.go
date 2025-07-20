package pg

import "context"

func (r *repo) ListenNotifications(ctx context.Context, channel string) (<-chan string, error) {
	conn, err := r.db.Acquire(ctx)
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
