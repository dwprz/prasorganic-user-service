ALTER TABLE users DROP CONSTRAINT users_pkey;

ALTER TABLE users ALTER COLUMN user_id TYPE INTEGER USING user_id::INTEGER;

ALTER TABLE users ADD CONSTRAINT users_pkey PRIMARY KEY (user_id);