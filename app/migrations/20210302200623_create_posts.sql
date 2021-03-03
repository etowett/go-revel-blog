-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table posts (
  id bigserial primary key,
  user_id bigint not null references users(id),
  category_id bigint not null references categories(id),
  title varchar(100) not null,
  content text not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz
);

create index posts_user_id_idx ON posts(user_id);
create index posts_category_id_idx ON posts(category_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists posts_category_id_idx;
drop index if exists posts_user_id_idx;

drop table if exists posts;
