-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tasks_comments (
    task_id BIGINT NOT NULL,
    comment TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,

    INDEX(task_id),

    CONSTRAINT fk_comments_task_id
    FOREIGN KEY (task_id) 
    REFERENCES tasks(task_id) 
    ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS tasks_comments;
-- +goose StatementEnd
