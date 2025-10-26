CREATE TABLE IF NOT EXISTS payments (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

	booking_id UUID NOT NULL,
	amount BIGINT NOT NULL,

	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (booking_id) REFERENCES bookings(id)
);

CREATE INDEX IF NOT EXISTS idx_payments_booking_id ON payments(booking_id);