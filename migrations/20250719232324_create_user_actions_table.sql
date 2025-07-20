-- +goose Up
-- +goose StatementBegin
create table if not exists user_actions (
  id bigserial primary key,
  user_id integer,
  action_type smallint,
  timestamp timestamptz not null,
  details jsonb not null
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table "user_actions";
-- +goose StatementEnd
