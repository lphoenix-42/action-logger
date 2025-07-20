package pg

import "context"

func (r *repo) ListenNotifications(ctx context.Context, channel string) (<-chan string, error) {
	return r.notifier.Listen(ctx, channel)
}
