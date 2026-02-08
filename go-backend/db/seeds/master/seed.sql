-- user statuses master data
insert into user_statuses (code, display_name, description, sort_order) values
  ('ACTIVE', 'Active', '正常に利用可能', 1),
  ('FROZEN', 'Frozen', '強制停止、ログイン不可', 2),
  ('DELETED', 'Deleted', '論理削除', 3) ON CONFLICT DO NOTHING;
