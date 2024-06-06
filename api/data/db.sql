CREATE TABLE IF NOT EXISTS users (
	b64 TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	name 	TEXT,
	password 	BLOB,
	gender 		TEXT,
	age 		INTEGER DEFAULT 0,
	first_name 	TEXT,
	last_name 	TEXT
);

CREATE TABLE IF NOT EXISTS posts (
	id 		INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id 		TEXT REFERENCES users(b64),
	content 	TEXT,
	date DATE
);

-- CREATE TABLE IF NOT EXISTS messages (
-- 	message_id INTEGER PRIMARY KEY AUTOINCREMENT,
-- 	sender_id TEXT REFERENCES users(b64),
-- 	receiver_id TEXT REFERENCES users(b64),
-- 	content TEXT,
-- 	time_stamp DATE
-- )
