-- name: AttachTagToBrag :exec
INSERT INTO brag_tags (brag_id, tag_id)
VALUES (?, ?);

-- name: DetachTagFromBrag :exec
DELETE FROM brag_tags WHERE brag_id = ? AND tag_id = ?;

-- name: DetachAllTagsFromBrag :exec
DELETE FROM brag_tags WHERE brag_id = ?;

-- name: DetachTagFromAllBrags :exec
DELETE FROM brag_tags WHERE tag_id = ?;

-- name: CountBragsByTag :one
SELECT COUNT(*) FROM brag_tags WHERE tag_id = ?;
