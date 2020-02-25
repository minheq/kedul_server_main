CREATE TABLE kedul_user (
  id UUID NOT NULL,
  full_name TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  country_code TEXT NOT NULL,
  profile_image_id TEXT NOT NULL DEFAULT '',
  is_phone_number_verified BOOLEAN NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  updated_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT "PK_kedul_user_1" PRIMARY KEY (id),
  CONSTRAINT "UN_kedul_user_1" UNIQUE (phone_number)
);