-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS team_members (
    user_id BIGINT NOT NULL,
    team_id BIGINT NOT NULL,
    joined_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, team_id),
    
    INDEX idx_team_users (team_id),

    CONSTRAINT fk_team_member_user_id 
    FOREIGN KEY (user_id) 
    REFERENCES users(user_id) 
    ON DELETE CASCADE,
    CONSTRAINT fk_team_member_team_id
    FOREIGN KEY (team_id) 
    REFERENCES teams(team_id) 
    ON DELETE CASCADE
) ENGINE=InnoDB;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS team_members;
-- +goose StatementEnd
