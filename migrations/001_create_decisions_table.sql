-- Creating the table
CREATE TABLE IF NOT EXISTS decisions (
    id BIGINT AUTO_INCREMENT PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    target_id VARCHAR(255) NOT NULL,
    decision ENUM('LIKE', 'PASS') NOT NULL,
    timestamp DATETIME NOT NULL,
    UNIQUE KEY unique_decision (user_id, target_id)
);