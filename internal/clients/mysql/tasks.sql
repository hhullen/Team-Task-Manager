-- name: AddNewTask :execresult
INSERT INTO tasks (assignee_id, created_by, team_id, subject, description, status)
VALUES (?, ?, ?, ?, ?, ?);


-- name: IsTeamMember :one
SELECT EXISTS(
    SELECT 1
    FROM team_members
    WHERE user_id = ? AND team_id = ?
);

-- name: GetTasksOfTeam :many
SELECT task_id, 
    assignee_id,
    created_by,
    team_id,
    subject,
    status,
    description,
    created_at
FROM tasks
WHERE team_id = ?;

-- name: GetTasks :many
SELECT task_id, 
    assignee_id,
    created_by,
    team_id,
    subject,
    status,
    description,
    created_at
FROM tasks
WHERE team_id = sqlc.arg(team_id)
  AND (status = sqlc.narg(status) OR sqlc.narg(status) IS NULL)
  AND (assignee_id = sqlc.narg(assignee_id) OR sqlc.narg(assignee_id) IS NULL)
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

