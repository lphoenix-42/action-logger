package pg

import (
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repo struct {
	db    *pgxpool.Pool
	query sq.StatementBuilderType
}

func New(db *pgxpool.Pool) *repo {
	return &repo{
		db:    db,
		query: sq.StatementBuilder.PlaceholderFormat(sq.Dollar),
	}
}
