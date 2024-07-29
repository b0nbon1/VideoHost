-- name: GetAllVideos :many
SELECT * FROM videos ORDER BY id;

-- name: CreateVideo :one
INSERT INTO videos (name, description, video_url) VALUES ($1, $2, $3) RETURNING *;

-- name: GetVideoById :one
SELECT * FROM videos WHERE id = $1 LIMIT 1;

-- name: DeleteVideoById :exec
DELETE FROM videos WHERE id = $1;

-- name: UpdateVideo :one
UPDATE videos SET name = $2, description = $3, video_url= $3 WHERE id = $1 RETURNING *;
