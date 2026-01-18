create table users (
  id uuid primary key,
  email varchar(256) unique not null,
  password_hash bytea not null,
  name varchar(256) not null,
  created_at timestamp not null default now()
);
