CREATE TABLE verification_code (
  id UUID NOT NULL,
  user_id UUID NOT NULL,
  code TEXT NOT NULL,
  verification_id TEXT NOT NULL,
  code_type TEXT NOT NULL,
  phone_number TEXT NOT NULL,
  country_code TEXT NOT NULL,
  created_at TIMESTAMPTZ NOT NULL,
  expired_at TIMESTAMPTZ NOT NULL,
  CONSTRAINT "PK_verification_code_1" PRIMARY KEY (id),
  CONSTRAINT "UN_verification_code_1" UNIQUE (phone_number)
);