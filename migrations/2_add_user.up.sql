INSERT INTO users (id, username, password_hash)
VALUES (1, 'testUser', 'QWERTYU')
ON CONFLICT DO NOTHING;

INSERT INTO boards (id, name)
VALUES (1, 'testBoard')
ON CONFLICT DO NOTHING;

INSERT INTO statuses (id, type)
VALUES (1, 'testStatus')
ON CONFLICT DO NOTHING;

INSERT INTO tasks (id, title, description, board_id, status_id, user_id)
VALUES (1, 'testTask', 'testDescription', 1, 1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO boards_users (id, user_id, board_id)
VALUES (1, 1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO user_token (id, user_id, access_token)
VALUES (1, 1, 'TESTTOKEN')
ON CONFLICT DO NOTHING;