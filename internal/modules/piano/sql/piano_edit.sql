-- name: CreatePianoEdit :exec
INSERT INTO piano_edits (id, piano_id, editor_user_id, operation, changes, summary)
VALUES (
    sqlc.arg(id)::ulid,
    sqlc.arg(piano_id)::ulid,
    sqlc.narg(editor_user_id)::ulid,
    sqlc.arg(operation)::piano_edit_operation,
    sqlc.narg(changes),
    sqlc.narg(summary)
);
