CREATE TABLE location (
  id UUID NOT NULL,
  business_id UUID NOT NULL,
  name TEXT NOT NULL,
  profile_image_id TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT "PK_location_1" PRIMARY KEY (id)
);