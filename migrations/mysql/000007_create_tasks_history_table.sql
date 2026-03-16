-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks_history (
    task_id BIGINT,
    changed_by BIGINT NOT NULL,
    payload JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    INDEX(task_id),

    CONSTRAINT fk_history_task_id
    FOREIGN KEY (task_id) 
    REFERENCES tasks(task_id) 
    ON DELETE CASCADE,
    CONSTRAINT fk_task_changed_by
    FOREIGN KEY (changed_by) 
    REFERENCES users(user_id) 
    ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks_history;
-- +goose StatementEnd
