CREATE TABLE company_invite (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  company_persist_key BIGINT NOT NULL REFERENCES company(persist_key),
  verify_token VARCHAR(255) NOT NULL,
  role_label VARCHAR(255) NOT NULL REFERENCES role(label),
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expire_date TIMESTAMP WITH TIME ZONE NOT NULL,
  UNIQUE (company_persist_key, verify_token)
);
