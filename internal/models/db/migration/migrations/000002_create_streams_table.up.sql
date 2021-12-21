CREATE TABLE IF NOT EXISTS streams(
	id bigserial PRIMARY KEY,
	name text,
	slug text,
	description text,
	start_date timestamptz,
	end_date timestamptz,
	manifest_path text,
	created_by int8,
	created_date timestamptz,
	updated_by int8,
	updated_date timestamptz
);


