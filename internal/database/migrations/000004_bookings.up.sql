CREATE TABLE IF NOT EXISTS bookings (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

	user_id UUID NOT NULL,
	field_id UUID NOT NULL,

	start_time TIMESTAMPTZ NOT NULL,
	end_time TIMESTAMPTZ NOT NULL,

	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (user_id) REFERENCES users(id),
    FOREIGN KEY (field_id) REFERENCES fields(id)
);

CREATE INDEX IF NOT EXISTS idx_bookings_user_id ON bookings(user_id);
CREATE INDEX IF NOT EXISTS idx_bookings_field_id ON bookings(field_id);