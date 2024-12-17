-- Creating indexes
-- Adding indexes on target_id (for the recipient) and user_id (for the actor) will drastically speed up lookups.
CREATE INDEX idx_target_id ON decisions(target_id);
-- CREATE INDEX idx_user_id ON decisions(user_id);
-- CREATE INDEX idx_target_decision ON decisions(target_id, decision);
-- CREATE INDEX idx_user_decision ON decisions(user_id, decision);