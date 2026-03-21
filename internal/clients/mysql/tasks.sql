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
    created_at,
    version
FROM tasks
WHERE team_id = sqlc.arg(team_id)
  AND (status = sqlc.narg(status) OR sqlc.narg(status) IS NULL)
  AND (assignee_id = sqlc.narg(assignee_id) OR sqlc.narg(assignee_id) IS NULL)
ORDER BY created_at DESC
LIMIT ? OFFSET ?;

-- name: AddChangeToTaskHistory :execresult
INSERT INTO tasks_history (task_id, changed_by, payload)
VALUES (?, ?, ?);

-- name: GetTaskForUpdate :one
SELECT assignee_id, team_id, subject, status, description, version
FROM tasks 
WHERE task_id = ? 
FOR UPDATE;

-- name: UpdateTask :execresult
UPDATE tasks 
SET assignee_id = ?,
    team_id = ?,
    subject = ?,
    status = ?,
    description = ?,
    version = version + 1,
    updated_at = NOW()
WHERE task_id = ?;

-- name: GetTask :one
SELECT assignee_id, 
    created_by,
    team_id,
    subject,
    status,
    description,
    created_at,
    version
FROM tasks 
WHERE task_id = ?
ORDER BY task_id;

-- name: GetTaskHistory :many
SELECT payload, changed_by, created_at
FROM tasks_history
WHERE task_id = ?
ORDER BY created_at
LIMIT ? OFFSET ?;

-- name: GetTasksComments :many
SELECT task_id, comment, created_at
FROM tasks_comments
WHERE task_id IN (sqlc.slice('task_ids'))
ORDER BY task_id;

-- name: AddTaskComment :execresult
INSERT INTO tasks_comments (task_id, comment)
VALUES (?, ?);
