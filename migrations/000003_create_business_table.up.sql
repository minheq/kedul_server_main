CREATE TABLE business (
  id UUID NOT NULL,
  user_id UUID NOT NULL,
  name TEXT NOT NULL,
  profile_image_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT "PK_business_1" PRIMARY KEY (id),
  CONSTRAINT "UN_business_1" UNIQUE (name)
);