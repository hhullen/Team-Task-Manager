-- name: AddNewTeam :execresult
INSERT INTO teams (owner_id, name, description)
VALUES (?, ?, ?);

-- name: AddMemberToTeam :execresult
INSERT INTO team_members (user_id, team_id)
VALUES (?, ?);

-- name: GetUserTeams :many
SELECT t.team_id, t.name, t.description
FROM team_members tm INNER JOIN teams t ON tm.team_id = t.team_id
WHERE tm.user_id = ?;

-- name: GetTeamOwner :one
select owner_id
FROM teams
WHERE team_id = ?;
