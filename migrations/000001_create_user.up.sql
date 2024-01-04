CREATE TABLE IF NOT EXISTS "user" (
	id serial PRIMARY KEY,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	name text NOT NULL,
	email text NOT NULL
);