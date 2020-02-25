CREATE TABLE employee_role (
  id UUID NOT NULL,
  location_id UUID NOT NULL,
  name TEXT NOT NULL,
  permission_ids TEXT [] NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT "PK_employee_role_1" PRIMARY KEY (id)
);