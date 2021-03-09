-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

create table comments (
  id bigserial primary key,
  user_id bigint not null references users(id),
  post_id bigint not null references posts(id),
  content text not null,
  created_at timestamptz not null default clock_timestamp(),
  updated_at timestamptz
);

create index comments_user_id_idx ON comments(user_id);
create index comments_post_id_idx ON comments(post_id);

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

drop index if exists comments_post_id_idx;
drop index if exists comments_user_id_idx;

drop table if exists comments;
