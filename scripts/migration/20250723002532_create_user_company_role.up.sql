CREATE TABLE user_company_role (
  persist_key BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
  user_persist_key BIGINT NOT NULL,
  company_persist_key BIGINT NOT NULL,
  role_label VARCHAR(255) NOT NULL,
  register_date TIMESTAMP WITH TIME ZONE NOT NULL,
  expire_date TIMESTAMP WITH TIME ZONE
);
