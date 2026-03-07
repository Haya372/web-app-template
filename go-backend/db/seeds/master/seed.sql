-- user statuses master data
insert into user_statuses (code, display_name, description, sort_order) values
  ('ACTIVE', 'Active', '正常に利用可能', 1),
  ('FROZEN', 'Frozen', '強制停止、ログイン不可', 2),
  ('DELETED', 'Deleted', '論理削除', 3) ON CONFLICT DO NOTHING;

-- roles master data
insert into roles (id, name, description) values
  ('00000000-0000-0000-0000-000000000001', 'admin', 'Full access to all resources'),
  ('00000000-0000-0000-0000-000000000002', 'viewer', 'Read-only access') ON CONFLICT DO NOTHING;

-- permissions master data
insert into permissions (id, code, description) values
  ('00000000-0000-0000-0001-000000000001', 'users:list', 'List users') ON CONFLICT DO NOTHING;

-- role_permissions: admin and viewer both get users:list
insert into role_permissions (role_id, permission_id) values
  ('00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0001-000000000001'),
  ('00000000-0000-0000-0000-000000000002', '00000000-0000-0000-0001-000000000001') ON CONFLICT DO NOTHING;
