package pg

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lphoenix-42/action-logger/internal/infrastructure/repository"
)

type repo struct {
	db       *pgxpool.Pool
	query    sq.StatementBuilderType
	notifier repository.Notifier
}

func New(db *pgxpool.Pool, notifier repository.Notifier) repository.ActionlogRepository {
	return &repo{
		db:       db,
		query:    sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
		notifier: notifier,
	}
}
