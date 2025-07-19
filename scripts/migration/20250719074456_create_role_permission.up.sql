CREATE TABLE role_permission (
  role_label VARCHAR(255) PRIMARY KEY,
  self_edit BOOLEAN NOT NULL,
  company_access BOOLEAN NOT NULL,
  company_invite BOOLEAN NOT NULL,
  company_edit BOOLEAN NOT NULL,
  priority INTEGER NOT NULL
);
