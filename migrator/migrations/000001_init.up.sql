CREATE TABLE IF NOT EXISTS holds (
                                     hold_id BIGINT NOT NULL,
                                     amount BIGINT NOT NULL,
                                     user_id VARCHAR NOT NULL,
                                     status VARCHAR NOT NULL,
                                     created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (hold_id)
);

CREATE TABLE IF NOT EXISTS users (
                                     user_id VARCHAR NOT NULL,
                                     balance BIGINT NOT NULL,
                                     username VARCHAR NOT NULL,
                                     password VARCHAR NOT NULL,
                                     created_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
                                     PRIMARY KEY (user_id)
);

CREATE INDEX IF NOT EXISTS idx_user_id ON users (user_id);
CREATE INDEX IF NOT EXISTS idx_hold_id ON holds (hold_id);