CREATE TABLE user_company_role (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_persist_key BIGINT NOT NULL REFERENCES users(persist_key),
  company_persist_key BIGINT NOT NULL REFERENCES company(persist_key),
  role_label VARCHAR(255) NOT NULL REFERENCES role(label),
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expire_date TIMESTAMP WITH TIME ZONE,
  UNIQUE (user_persist_key, company_persist_key, role_label)
);
