// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: videos.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createVideo = `-- name: CreateVideo :one
INSERT INTO videos (name, description, video_url) VALUES ($1, $2, $3) RETURNING id, name, description, video_url, createdat, updatedat
`

type CreateVideoParams struct {
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
	VideoUrl    string         `json:"video_url"`
}

func (q *Queries) CreateVideo(ctx context.Context, arg CreateVideoParams) (Video, error) {
	row := q.queryRow(ctx, q.createVideoStmt, createVideo, arg.Name, arg.Description, arg.VideoUrl)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.VideoUrl,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const deleteVideoById = `-- name: DeleteVideoById :exec
DELETE FROM videos WHERE id = $1
`

func (q *Queries) DeleteVideoById(ctx context.Context, id uuid.UUID) error {
	_, err := q.exec(ctx, q.deleteVideoByIdStmt, deleteVideoById, id)
	return err
}

const getAllVideos = `-- name: GetAllVideos :many
SELECT id, name, description, video_url, createdat, updatedat FROM videos ORDER BY id
`

func (q *Queries) GetAllVideos(ctx context.Context) ([]Video, error) {
	rows, err := q.query(ctx, q.getAllVideosStmt, getAllVideos)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Video
	for rows.Next() {
		var i Video
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Description,
			&i.VideoUrl,
			&i.Createdat,
			&i.Updatedat,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getVideoById = `-- name: GetVideoById :one
SELECT id, name, description, video_url, createdat, updatedat FROM videos WHERE id = $1 LIMIT 1
`

func (q *Queries) GetVideoById(ctx context.Context, id uuid.UUID) (Video, error) {
	row := q.queryRow(ctx, q.getVideoByIdStmt, getVideoById, id)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.VideoUrl,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}

const updateVideo = `-- name: UpdateVideo :one
UPDATE videos SET name = $2, description = $3, video_url= $3 WHERE id = $1 RETURNING id, name, description, video_url, createdat, updatedat
`

type UpdateVideoParams struct {
	ID          uuid.UUID      `json:"id"`
	Name        string         `json:"name"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) UpdateVideo(ctx context.Context, arg UpdateVideoParams) (Video, error) {
	row := q.queryRow(ctx, q.updateVideoStmt, updateVideo, arg.ID, arg.Name, arg.Description)
	var i Video
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Description,
		&i.VideoUrl,
		&i.Createdat,
		&i.Updatedat,
	)
	return i, err
}
