DELETE FROM task_types WHERE task_type IN (
  'subscribe_to_tg',
  'subscribe_to_x',
  'subscribe_to_yt'
);
DELETE FROM users WHERE id IN (
  'bc198810-5608-4ad2-b251-b0e87843d0ae',
  'dd74ba03-96ad-4632-a56c-d49a658c17e1'
);