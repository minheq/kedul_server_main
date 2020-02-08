CREATE TABLE verification_code (
	id							UUID												NOT NULL,
	account_id			UUID												NOT NULL,
	code						TEXT												NOT NULL,
	verification_id	TEXT												NOT NULL,
	code_type				TEXT												NOT NULL,
	phone_number		TEXT												NOT NULL,
	country_code		TEXT												NOT NULL,
	created_at			TIMESTAMP WITH TIME ZONE		NOT NULL,
	expired_at			TIMESTAMP WITH TIME ZONE		NOT NULL,	

	CONSTRAINT	"PK_verification_code_1"	PRIMARY KEY	(id),
	CONSTRAINT	"UN_verification_code_1"	UNIQUE	(phone_number)
);