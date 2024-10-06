INSERT INTO statuses (id, type)
VALUES (1, 'in process')
ON CONFLICT DO NOTHING;

INSERT INTO statuses (id, type)
VALUES (2, 'done')
ON CONFLICT DO NOTHING;