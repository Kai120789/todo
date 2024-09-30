INSERT INTO users (id, username, password_hash)
VALUES (1, 'testUser', 'QWERTYU')
ON CONFLICT DO NOTHING;

INSERT INTO boards (id, name, user_id)
VALUES (1, 'testBoard', 1)
ON CONFLICT DO NOTHING;

INSERT INTO statuses (id, type)
VALUES (1, 'testStatus')
ON CONFLICT DO NOTHING;

INSERT INTO tasks (id, title, description, board_id, status_id, user_id)
VALUES (1, 'testTask', 'testDescription', 1, 1, 1)
ON CONFLICT DO NOTHING;