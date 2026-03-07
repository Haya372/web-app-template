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

create table roles (
  id uuid primary key,
  name varchar(64) not null unique,
  description text,
  created_at timestamp not null default now(),
  updated_at timestamp not null default now()
);

create table permissions (
  id uuid primary key,
  code varchar(128) not null unique,
  description text,
  created_at timestamp not null default now()
);

create table role_permissions (
  role_id uuid not null references roles(id),
  permission_id uuid not null references permissions(id),
  primary key (role_id, permission_id)
);

create table user_roles (
  user_id uuid not null references users(id),
  role_id uuid not null references roles(id),
  primary key (user_id, role_id)
);
