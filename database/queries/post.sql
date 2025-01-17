-- name: CreatePost :one
INSERT INTO post (id, title, description)
VALUES ($1, $2, $3)
RETURNING id, title, description, created_at, updated_at;

-- name: GetPost :one
SELECT id, title, description, created_at, updated_at
FROM post
WHERE id = $1;

-- name: ListPosts :many
SELECT id, title, description, created_at, updated_at
FROM post
ORDER BY created_at DESC;

-- name: UpdatePost :one
UPDATE post SET
  title = $2,
  description = $3,
  updated_at = NOW()
WHERE id = $1
RETURNING id, title, description, created_at, updated_at;

-- name: DeletePost :execrows
DELETE FROM post WHERE id = $1;
