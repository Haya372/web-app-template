create table user_statuses (
  code varchar(32) primary key,
  display_name varchar(64) not null,
  description text,
  sort_order smallint not null default 0,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

insert into user_statuses (code, display_name, description, sort_order) values
  ('ACTIVE', 'Active', '正常に利用可能', 1),
  ('FROZEN', 'Frozen', '強制停止、ログイン不可', 2),
  ('DELETED', 'Deleted', '論理削除', 3);

create table users (
  id uuid primary key,
  email varchar(256) unique not null,
  password_hash bytea not null,
  name varchar(256) not null,
  status_code varchar(32) not null default 'ACTIVE' references user_statuses(code),
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);
