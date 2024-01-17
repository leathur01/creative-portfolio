CREATE TABLE IF NOT EXISTS portfolio (
	id serial PRIMARY KEY,
	name text NOT NULL,
	created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
	user_id serial NOT NULL, 
	CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES "user"(id)
)