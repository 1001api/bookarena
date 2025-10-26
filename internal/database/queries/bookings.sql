-- name: CreateBooking :one
INSERT INTO bookings (
    user_id, 
    field_id, 
    start_time, 
    end_time,
    total_price
)
VALUES (
    @user_id, 
    @field_id, 
    @start_time, 
    @end_time,
    @total_price
)
RETURNING id;

-- name: GetBookingByID :one
SELECT 
    b.id,
    b.user_id,
    b.field_id,
    b.start_time,
    b.end_time,
    b.total_price,
    b.created_at,
    b.updated_at,
    b.deleted_at,
    f.name AS field_name,
    f.type AS field_type,
    f.location AS field_location,
    f.location_lat AS field_location_lat,
    f.location_lon AS field_location_lon,
    f.price_per_hour AS field_price_per_hour,
    u.username AS user_username,
    u.email_enc::text AS user_email
FROM bookings b
LEFT JOIN fields f ON b.field_id = f.id
LEFT JOIN users u ON b.user_id = u.id
WHERE b.id = @id::uuid AND b.deleted_at IS NULL;

-- name: ListBookings :many
SELECT 
    b.id,
    b.user_id,
    b.field_id,
    b.start_time,
    b.end_time,
    b.total_price,
    b.created_at,
    b.updated_at,
    b.deleted_at,
    f.name AS field_name,
    f.type AS field_type,
    f.location AS field_location,
    f.location_lat AS field_location_lat,
    f.location_lon AS field_location_lon,
    f.price_per_hour AS field_price_per_hour,
    u.username AS user_username,
    u.email_enc::text AS user_email
FROM bookings b
LEFT JOIN fields f ON b.field_id = f.id
LEFT JOIN users u ON b.user_id = u.id
WHERE b.deleted_at IS NULL
ORDER BY b.created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: ListBookingsByUser :many
SELECT 
    b.id,
    b.user_id,
    b.field_id,
    b.start_time,
    b.end_time,
    b.total_price,
    f.name AS field_name,
    f.type AS field_type,
    f.location AS field_location,
    f.location_lat AS field_location_lat,
    f.location_lon AS field_location_lon,
    f.price_per_hour AS field_price_per_hour,
    b.created_at,
    b.updated_at
FROM bookings b
LEFT JOIN fields f ON b.field_id = f.id
WHERE b.user_id = @user_id::uuid AND b.deleted_at IS NULL
ORDER BY b.created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: ListBookingsByField :many
SELECT 
    b.id,
    b.user_id,
    b.field_id,
    b.start_time,
    b.end_time,
    b.total_price,
    b.created_at,
    b.updated_at,
    b.deleted_at,
    f.name AS field_name,
    f.type AS field_type,
    f.location AS field_location,
    f.location_lat AS field_location_lat,
    f.location_lon AS field_location_lon,
    f.price_per_hour AS field_price_per_hour,
    u.username AS user_username,
    u.email_enc::text AS user_email
FROM bookings b
LEFT JOIN fields f ON b.field_id = f.id
LEFT JOIN users u ON b.user_id = u.id
WHERE b.field_id = @field_uuid::uuid AND b.deleted_at IS NULL
ORDER BY b.created_at DESC
LIMIT @limit_count OFFSET @offset_count;

-- name: DeleteBooking :one
UPDATE bookings
SET deleted_at = NOW()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id;

-- name: CheckBookingConflict :one
SELECT EXISTS (
    SELECT 1 FROM bookings
    WHERE field_id = @field_id::uuid
      AND deleted_at IS NULL
      AND tstzrange(start_time, end_time) && tstzrange(@start_time::timestamptz, @end_time::timestamptz)
) AS booking_conflict;