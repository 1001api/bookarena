CREATE TABLE IF NOT EXISTS payments (
	id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

	booking_id UUID NOT NULL,
	user_id UUID NOT NULL,
	total_price BIGINT NOT NULL,
	status VARCHAR(255) NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'paid', 'failed')),

	created_at TIMESTAMPTZ DEFAULT NOW(),
	updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ,

    FOREIGN KEY (booking_id) REFERENCES bookings(id),
    FOREIGN KEY (user_id) REFERENCES users(id)
);

CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments(user_id);
CREATE INDEX IF NOT EXISTS idx_payments_booking_id ON payments(booking_id);
