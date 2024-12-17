-- Inserting data
INSERT IGNORE INTO decisions (user_id, target_id, decision, timestamp) 
VALUES 
    ('user1', 'target1', 'LIKE', NOW()),
    ('user2', 'target2', 'PASS', NOW()),
    ('user3', 'target3', 'LIKE', NOW());