-- name: GetAnnotation :many
SELECT *
FROM annotations
WHERE id = ?
LIMIT 1;
-- name: GetAnnotations :exec
SELECT *
FROM annotations
ORDER BY keyword;
-- name: CreateAnnotation :one
INSERT INTO annotations (
        filename,
        keyword,
        content,
        notes
    )
VALUES (?, ?, ?, ?)
RETURNING *;
-- name: UpdateAnnotation :exec
UPDATE annotations
set filename = ?,
    keyword = ?,
    content = ?,
    notes = ?,
    updated_at = CURRENT_TIMESTAMP
WHERE id = ?;
 -- name: DeleteAnnotation :exec
DELETE FROM annotations
WHERE id = ?
RETURNING *;