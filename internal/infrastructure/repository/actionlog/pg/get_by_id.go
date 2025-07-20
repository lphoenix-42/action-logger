package pg

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/converter"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/repository/actionlog/pg/schema"
	"github.com/lphoenix-42/action-logger/internal/service/model"
)

func (r *repo) GetByID(ctx context.Context, id int64) (*model.Action, error) {
	action, err := r.getFromDbByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return action, nil
}

func (r *repo) getFromDbByID(ctx context.Context, id int64) (*model.Action, error) {
	sqlBuilder := r.query.
		Select("id", "user_id", "action_type", "timestamp", "details").
		From(`"user_actions"`).
		Where(sq.Eq{"id": id})

	query, args, err := sqlBuilder.ToSql()
	if err != nil {
		return nil, err
	}

	var actionSchema schema.ActionSchema
	err = pgxscan.Get(ctx, r.db, &actionSchema, query, args...)
	if err != nil {
		return nil, err
	}

	action, err := converter.ActionFromSchemaToModel(&actionSchema)
	if err != nil {
		return nil, err
	}

	return action, nil
}
