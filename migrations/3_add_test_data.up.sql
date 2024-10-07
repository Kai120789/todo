-- test inserts in all tables

INSERT INTO users (id, username, tg_name, password_hash)
VALUES (1, 'testUser', 'test', 'QWERTYU')
ON CONFLICT DO NOTHING;

INSERT INTO boards (id, name)
VALUES (1, 'testBoard')
ON CONFLICT DO NOTHING;

INSERT INTO tasks (id, title, description, board_id, status_id, user_id)
VALUES (1, 'testTask', 'testDescription', 1, 1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO boards_users (id, user_id, board_id)
VALUES (1, 1, 1)
ON CONFLICT DO NOTHING;

INSERT INTO user_token (id, user_id, refresh_token)
VALUES (1, 1, 'header.payload.secret')
ON CONFLICT DO NOTHING;