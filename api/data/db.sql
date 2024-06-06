CREATE TABLE IF NOT EXISTS users (
	b64 TEXT PRIMARY KEY,
	email TEXT NOT NULL UNIQUE,
	username 	TEXT,
	password 	TEXT,
	gender 		TEXT,
	age 		INTEGER DEFAULT 0,
	firstname 	TEXT,
	lastname 	TEXT
);

CREATE TABLE IF NOT EXISTS posts (
	uuid 		TEXT PRIMARY KEY,
	userid 		TEXT REFERENCES users(b64),
	content 	TEXT,
	date DATE
);

