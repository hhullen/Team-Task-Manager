-- name: GetTeamTasksCompletedFor7Days :many
SELECT t.name,
	   COALESCE(tm.members, 0) AS members,
	   COALESCE(ts.completed_in_7_days, 0) AS completed_in_7_days
FROM teams t
	LEFT JOIN
		(SELECT team_id, COUNT(*) AS members
		 FROM team_members
		 GROUP BY team_id) tm ON t.team_id = tm.team_id
	LEFT JOIN 
		(SELECT team_id, COUNT(*) completed_in_7_days
		 FROM tasks
		 WHERE status = "completed" 
		 	AND updated_at >= NOW() - INTERVAL 7 DAY
		 GROUP BY team_id) ts ON t.team_id = ts.team_id;

-- nmae: GetTop3ByCreatedTeamTasksLastMonth :many
WITH user_stats AS (
	SELECT tm.team_id, tm.user_id AS created_by, COALESCE(COUNT(task_id), 0) as total_tasks
	FROM team_members tm
		LEFT JOIN tasks t ON t.created_by = tm.user_id
			AND t.team_id = tm.team_id
			AND t.created_at >= NOW() - INTERVAL 1 MONTH
	GROUP BY tm.team_id, tm.user_id
),
user_ranked AS (
	SELECT team_id, 
		   created_by,
		   total_tasks,
		   DENSE_RANK() OVER (
			   PARTITION BY team_id
			   ORDER BY total_tasks DESC
		   ) AS ranks
	FROM user_stats)
SELECT 
	ur.team_id,
	t.name AS team_name,
	ur.created_by AS user_id,
	ur.total_tasks,
	ur.ranks
FROM user_ranked ur
JOIN teams t ON t.team_id = ur.team_id 
WHERE ranks <= 3
ORDER BY ur.team_id, ur.ranks;

-- name: GetTasksAssigneeNotTeamMember :many
SELECT t.*
FROM tasks t
	LEFT JOIN team_members tm ON t.team_id = tm.team_id AND t.assignee_id = tm.user_id
WHERE tm.user_id IS NULL;
