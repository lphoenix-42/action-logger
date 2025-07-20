package pg

import (
	"context"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/converter"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/repository/actionlog/pg/schema"
	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (r *repo) QueryActions(ctx context.Context, filter *model.ActionSearch) ([]*model.Action, error) {
	query, args, err := r.buildQuery(filter)
	if err != nil {
		return nil, err
	}

	var rows []schema.ActionSchema
	if err := pgxscan.Select(ctx, r.db, &rows, query, args...); err != nil {
		return nil, err
	}

	var actions []*model.Action
	for _, row := range rows {
		action, err := converter.ActionFromSchemaToModel(&row)
		if err != nil {
			return nil, err
		}
		actions = append(actions, action)
	}
	return actions, nil
}

func (r *repo) buildQuery(filter *model.ActionSearch) (string, []any, error) {
	q := r.query.
		Select("id", "user_id", "action_type", "timestamp", "details").
		From(`"user_actions"`).
		OrderBy("timestamp DESC")

	if filter.UserID != nil {
		q = q.Where(sq.Eq{"user_id": *filter.UserID})
	}
	if len(filter.ActionTypes) > 0 {
		q = q.Where(sq.Eq{"action_type": filter.ActionTypes})
	}
	if filter.From != nil {
		q = q.Where(sq.GtOrEq{"timestamp": *filter.From})
	}
	if filter.To != nil {
		q = q.Where(sq.LtOrEq{"timestamp": *filter.To})
	}

	for _, f := range filter.DetailsFilters {
		if len(f.Path) == 0 {
			continue
		}
		path := "{" + strings.Join(f.Path, ",") + "}"
		q = q.Where(sq.Expr("details #>> ? = ?", path, f.Value))
	}

	return q.PlaceholderFormat(sq.Dollar).ToSql()
}
