-- name: CreateField :one
INSERT INTO fields (
    image_url, 
    name, 
    type, 
    description, 
    location,
    location_lat, 
    location_lon, 
    price_per_hour
)
VALUES (
    @image_url,
    @name,
    @type,
    @description,
    @location,
    @location_lat,
    @location_lon,
    @price_per_hour
) RETURNING id;

-- name: GetFieldByID :one
SELECT 
    id,
    image_url,
    name,
    type,
    description,
    location,
    location_lat,
    location_lon,
    price_per_hour,
    created_at,
    updated_at
FROM fields
WHERE id = @id::uuid AND deleted_at IS NULL;

-- name: ListFields :many
SELECT 
    id,
    image_url,
    name,
    type,
    description,
    price_per_hour,
    created_at,
    updated_at
FROM fields
WHERE deleted_at IS NULL
ORDER BY created_at DESC
OFFSET @offset_count LIMIT @limit_count;

-- name: UpdateField :one
UPDATE fields
SET
    image_url = COALESCE(@image_url, image_url),
    name = COALESCE(@name, name),
    type = COALESCE(@type, type),
    description = COALESCE(@description, description),
    location = COALESCE(@location, location),
    location_lat = COALESCE(@location_lat, location_lat),
    location_lon = COALESCE(@location_lon, location_lon),
    price_per_hour = COALESCE(@price_per_hour, price_per_hour),
    updated_at = NOW()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id;

-- name: DeleteField :one
UPDATE fields
SET deleted_at = NOW()
WHERE id = @id::uuid AND deleted_at IS NULL
RETURNING id;