CREATE TABLE user_password (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_persist_key BIGINT NOT NULL REFERENCES users(persist_key),
  password VARCHAR(255) NOT NULL,
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expire_date TIMESTAMP WITH TIME ZONE,
  unique (user_persist_key, password)
);
