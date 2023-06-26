CREATE TABLE IF NOT EXISTS {{ index .Options "Namespace" }}.legacy_credentials (
	id uuid NOT NULL DEFAULT gen_random_uuid(),
	email varchar(255) NULL,
	phone VARCHAR(15) NULL DEFAULT NULL,
	encrypted_password varchar(255) NULL,
	created_at timestamptz NULL DEFAULT NOW(),
	updated_at timestamptz NULL,
	CONSTRAINT legacy_credentials_pkey PRIMARY KEY (id)
);
comment on table {{ index .Options "Namespace" }}.legacy_credentials is 'Auth: Stores legacy user login data within a secure schema.';
