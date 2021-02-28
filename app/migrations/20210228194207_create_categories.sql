-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table categories (
  id bigserial primary key,
  user_id bigint not null references users(id),
  name varchar(100) not null,
  description text not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz
);

create index categories_user_id_idx ON categories(user_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists categories_user_id_idx;

drop table if exists categories;
