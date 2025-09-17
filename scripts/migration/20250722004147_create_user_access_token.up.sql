CREATE TABLE user_access_token (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_persist_key BIGINT NOT NULL REFERENCES users(persist_key),
  access_token VARCHAR(4096) NOT NULL,
  source_update_date TIMESTAMP WITH TIME ZONE NOT NULL,
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expire_date TIMESTAMP WITH TIME ZONE NOT NULL,
  UNIQUE (user_persist_key, access_token)
);
