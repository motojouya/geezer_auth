CREATE TABLE user_refresh_token (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_persist_key BIGINT NOT NULL REFERENCES users(persist_key),
  refresh_token VARCHAR(255) NOT NULL,
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expire_date TIMESTAMP WITH TIME ZONE,
  unique (refresh_token)
);
