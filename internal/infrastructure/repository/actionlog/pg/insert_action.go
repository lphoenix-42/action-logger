package pg

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/lphoenix-42/action-logger/internal/service/model"
	"github.com/pkg/errors"
)

func (r *repo) InsertAction(ctx context.Context, info *model.ActionInfo) (int64, error) {
	detailsBytes, err := json.Marshal(info.Details)
	if err != nil {
		return 0, err
	}

	queryBuilder := r.query.
		Insert("user_actions").
		Columns("user_id", "action_type", "timestamp", "details").
		Values(info.UserID, info.ActionType, info.Timestamp, detailsBytes).
		Suffix("RETURNING id")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var id int64
	err = pgxscan.Get(ctx, r.db, &id, query, args...)
	if err != nil {
		return 0, err
	}

	if err := r.notifier.Notify(ctx, "new_action", fmt.Sprintf("%d", id)); err != nil {
		return 0, errors.Wrap(err, "inserted but failed to notify")
	}

	return id, nil
}
