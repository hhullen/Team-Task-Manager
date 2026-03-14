-- name: AddNewTeam :execresult
INSERT INTO teams (owner_id, name, description)
VALUES (?, ?, ?);

-- name: AddMemberToTeam :execresult
INSERT INTO team_members (user_id, team_id)
VALUES (?, ?);

-- name: GetUserTeams :many
SELECT t.team_id, t.name, t.description
FROM teams t
WHERE t.owner_id = ?
ORDER BY t.name;

-- name: GetTeamOwner :one
select owner_id
FROM teams
WHERE team_id = ?;
