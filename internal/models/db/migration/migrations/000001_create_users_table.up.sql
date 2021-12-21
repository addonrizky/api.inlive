CREATE TABLE IF NOT EXISTS users(
	id bigserial PRIMARY KEY,
	username text,
	password text,
	name text,
	login_type varchar(255),
	email text,
	role int8,
	is_active bool,
	is_login bool,
	register_date timestamptz,
	updated_date timestamptz
);