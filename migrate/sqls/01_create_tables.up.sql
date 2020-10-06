CREATE TABLE account (
    uid BIGINT NOT NULL,
    session_token VARCHAR(512) NOT NULL,
    tag VARCHAR(64) NOT NULL,
    PRIMARY KEY (uid, tag)
);

CREATE TABLE user (
    uid BIGINT NOT NULL,
    user_name VARCHAR(32) NOT NULL,
    is_block BOOLEAN NOT NULL,
    max_account TINYINT NOT NULL,
    n_account TINYINT NOT NULL,
    is_admin BOOLEAN NOT NULL,
    allow_polling BOOLEAN NOT NULL,
    PRIMARY KEY (uid)
);

CREATE TABLE runtime (
    uid BIGINT NOT NULL,
    session_token VARCHAR(512) NOT NULL,
    iksm CHARACTER(40) NOT NULL,
    language VARCHAR(10) NOT NULL,
    PRIMARY KEY (uid)
);

CREATE INDEX idx_user_name ON user(user_name);