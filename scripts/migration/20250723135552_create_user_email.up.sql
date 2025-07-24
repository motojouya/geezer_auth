CREATE TABLE user_email (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_persist_key BIGINT NOT NULL REFERENCES users(persist_key),
  email VARCHAR(255) NOT NULL,
  verify_token VARCHAR(255) NOT NULL,
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  verify_date TIMESTAMP WITH TIME ZONE,
  expire_date TIMESTAMP WITH TIME ZONE,
  UNIQUE (user_persist_key, email)
);
