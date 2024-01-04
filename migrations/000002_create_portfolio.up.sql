CREATE TABLE IF NOT EXISTS portfolio (
	id serial PRIMARY KEY,
	name text NOT NULL,
	user_id serial NOT NULL, 
	CONSTRAINT fk_user
      FOREIGN KEY(user_id) 
	  REFERENCES "user"(id)
)