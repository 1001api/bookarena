-- name: CreatePayment :one
INSERT INTO payments (
    booking_id,
    user_id,
    total_price,
    status
) VALUES (
    @booking_id,
    @user_id,
    @total_price,
    COALESCE(@status, 'pending')
)
RETURNING id;

-- name: GetPaymentByID :one
SELECT
    p.id,

    p.booking_id,
    p.user_id,
    p.total_price,
    p.status,

    b.start_time,
    b.end_time,

    f.id AS field_id,
    f.name AS field_name,
    f.type AS field_type,
    f.price_per_hour,

    u.username AS user_username,
    u.email_enc::text AS user_email,

    p.created_at AS payment_created_at,
    p.updated_at AS payment_updated_at
FROM payments p
JOIN bookings b ON b.id = p.booking_id
JOIN fields f ON f.id = b.field_id
JOIN users u ON u.id = p.user_id
WHERE p.id = @id::uuid AND p.deleted_at IS NULL;

-- name: ListPaymentsForAdmin :many
SELECT
    p.id,
    p.booking_id,
    p.user_id,
    p.total_price,
    p.status,
    u.email_enc::text AS user_email,
    b.start_time AS booking_start_time,
    b.end_time AS booking_end_time,
    p.created_at,
    p.updated_at
FROM payments p
JOIN users u ON u.id = p.user_id
JOIN bookings b ON b.id = p.booking_id
WHERE p.deleted_at IS NULL
ORDER BY p.created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: ListPaymentsForUser :many
SELECT
    p.id,
    p.booking_id,
    p.user_id,
    p.total_price,
    p.status,
    b.start_time AS booking_start_time,
    b.end_time AS booking_end_time,
    p.created_at,
    p.updated_at
FROM payments p
JOIN users u ON u.id = p.user_id
JOIN bookings b ON b.id = p.booking_id
WHERE p.deleted_at IS NULL AND p.user_id = @user_uuid::uuid
ORDER BY p.created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: UpdatePaymentStatus :one
UPDATE payments
SET 
    status = @status,
    updated_at = NOW()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id;

-- name: DeletePayment :one
UPDATE payments
SET deleted_at = NOW()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id;