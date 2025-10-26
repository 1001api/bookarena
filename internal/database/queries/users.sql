-- name: GetUsers :many
SELECT
    id,
    email_enc AS email,
    username,
    role AS role_name,
    last_login_at,
    created_at,
    last_updated_at AS updated_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
OFFSET @offset_count LIMIT @limit_count;

-- name: SearchUser :many
SELECT
    id,
    username,
    email_enc AS email,
    role AS role_name,
    last_login_at,
    created_at,
    last_updated_at AS updated_at
FROM users
WHERE deleted_at IS NULL
  AND (
      username ILIKE '%' || @query::text || '%'
  )
ORDER BY created_at DESC
OFFSET @offset_count LIMIT @limit_count;

-- name: CreateUser :one
INSERT INTO users (
    email_hash,
    email_enc,
    username,
    password_hash,
    role
)
VALUES (
    @email_hash::text,
    @email_enc,
    @username,
    @password_hash::text,
    @role::text
) RETURNING id;

-- name: UpdateUser :exec
UPDATE users
SET
    username = CASE
        WHEN @username::text IS NOT NULL
            AND @username::text != ''
            AND @username::text != username
        THEN @username::text
        ELSE username
    END,
    email_hash = CASE
        WHEN @email_hash::text IS NOT NULL
            AND @email_hash::text != ''
            AND @email_hash::text != email_hash
        THEN @email_hash::text
        ELSE email_hash
    END,
    email_enc = CASE
        WHEN @email_enc::bytea IS NOT NULL
            AND @email_enc::bytea != email_enc
        THEN @email_enc::bytea
        ELSE email_enc
    END,
    last_updated_at = CASE
        WHEN (
            (@username::text IS NOT NULL AND @username::text != '' AND @username::text != username)
            OR (@email_hash::text IS NOT NULL AND @email_hash::text != '' AND @email_hash::text != email_hash)
            OR (@email_enc::bytea IS NOT NULL AND @email_enc::bytea != email_enc)
        ) THEN CURRENT_TIMESTAMP
        ELSE last_updated_at
    END
WHERE id = @id;

-- name: DeleteUser :exec
UPDATE users
SET 
    deleted_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: RestoreUser :exec
UPDATE users
SET 
    deleted_at = NULL
WHERE id = @id;

-- name: CheckUserExists :one
SELECT EXISTS (
    SELECT 1
    FROM users u
    WHERE u.deleted_at IS NULL
      AND (
          (@email_hash::text IS NOT NULL AND @email_hash != '' AND email_hash = @email_hash)
          OR (@username::text IS NOT NULL AND @username != '' AND username = @username)
      )
) AS exists;

-- name: CheckUserExistsByID :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE id = @id::uuid AND deleted_at IS NULL
) AS exists;

-- name: GetUserByEmail :one
SELECT
    u.id,
    u.username,
    u.email_enc AS email,
    u.role AS role_name,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
WHERE u.email_hash = @email_hash
  AND u.deleted_at IS NULL
LIMIT 1;

-- name: GetUserByID :one
SELECT
    u.id,
    u.username,
    u.email_enc AS email,
    u.role AS role_name,
    u.last_login_at,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
WHERE u.id = @id
  AND u.deleted_at IS NULL;

-- name: GetUserByIdentifier :one
SELECT
    u.id,
    u.username,
    u.email_enc AS email,
    u.role AS role_name,
    u.password_hash,
    u.locked_until,
    u.failed_login_attempts,
    u.created_at,
    u.last_updated_at AS updated_at
FROM users u
WHERE u.deleted_at IS NULL
  AND (
      u.id = @id
      OR (NULLIF(@email_hash, '') IS NOT NULL AND u.email_hash = @email_hash)
      OR (NULLIF(@username, '') IS NOT NULL AND u.username = @username)
  )
LIMIT 1;

-- name: UpdateLastLogin :exec
UPDATE users
SET last_login_at = CURRENT_TIMESTAMP
WHERE id = @id;

-- name: IncrementFailedLoginCount :exec
UPDATE users
SET failed_login_attempts = failed_login_attempts + 1
WHERE id = @id;

-- name: ResetFailedLoginCount :exec
UPDATE users
SET 
    failed_login_attempts = 0,
    locked_until = NULL
WHERE id = @id;

-- name: LockUser :exec
UPDATE users
SET 
    locked_until = @locked_until::timestamptz,
    failed_login_attempts = @failed_attempts::int
WHERE id = @id;