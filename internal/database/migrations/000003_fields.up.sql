CREATE TABLE IF NOT EXISTS fields (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    image_url TEXT,
	name VARCHAR(100) NOT NULL,
	type VARCHAR(100) NOT NULL,
    description TEXT,
    location VARCHAR(255),
    location_lat DOUBLE PRECISION,
    location_lon DOUBLE PRECISION,
    price_per_hour BIGINT NOT NULL,

	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_fields_name ON fields(name);