-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks (
    task_id BIGINT AUTO_INCREMENT NOT NULL,
    assignee_id BIGINT NOT NULL,
    created_by BIGINT NOT NULL,
    team_id BIGINT NOT NULL,
    subject VARCHAR(255) NOT NULL,
    status VARCHAR(16) NOT NULL DEFAULT 'todo',
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 1,

    PRIMARY KEY (team_id, task_id),
    INDEX idx_team_task (task_id),
    INDEX idx_status_update_team (status, updated_at, team_id),

    CONSTRAINT fk_task_assignee_id
    FOREIGN KEY (assignee_id) 
    REFERENCES users(user_id) 
    ON DELETE CASCADE,
    CONSTRAINT fk_task_created_by
    FOREIGN KEY (created_by) 
    REFERENCES users(user_id) 
    ON DELETE CASCADE,
    CONSTRAINT fk_task_team_id
    FOREIGN KEY (team_id) 
    REFERENCES teams(team_id) 
    ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks;
-- +goose StatementEnd
