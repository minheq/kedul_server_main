CREATE TABLE kedul_user (
	id                          UUID                    		NOT NULL,
	full_name                   TEXT                    		NOT NULL,
	phone_number                TEXT                    		NOT NULL,
	country_code								TEXT												NOT NULL,
	is_phone_number_verified    BOOLEAN											NOT NULL,
	created_at									TIMESTAMP WITH TIME ZONE		NOT NULL,
	updated_at									TIMESTAMP WITH TIME ZONE		NOT NULL,	

	CONSTRAINT "PK_user_1" 	PRIMARY KEY (id),
	CONSTRAINT "UN_user_1" 	UNIQUE (phone_number)
);