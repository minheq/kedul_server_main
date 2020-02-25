CREATE TABLE employee (
  id UUID NOT NULL,
  location_id UUID NOT NULL,
  name TEXT NOT NULL,
  user_id UUID,
  employee_role_id UUID NOT NULL,
  profile_image_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT "PK_employee_1" PRIMARY KEY (id)
);