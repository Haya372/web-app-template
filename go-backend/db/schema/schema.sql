create table user_statuses (
  code varchar(32) primary key,
  display_name varchar(64) not null,
  description text,
  sort_order smallint not null default 0,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

create table users (
  id uuid primary key,
  email varchar(256) unique not null,
  password_hash bytea not null,
  name varchar(256) not null,
  status_code varchar(32) not null default 'ACTIVE' references user_statuses(code),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);
