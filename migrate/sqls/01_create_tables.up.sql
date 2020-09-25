CREATE TABLE account (
    uid BIGINT NOT NULL,
    iksm CHARACTER(40) NOT NULL,
    PRIMARY KEY (uid)
);

CREATE TABLE username (
    uid BIGINT NOT NULL,
    iksm VARCHAR(32) NOT NULL,
    PRIMARY KEY (uid)
);

CREATE TABLE status (
    uid BIGINT NOT NULL,
    iksm VARCHAR(32) NOT NULL,
    is_block BOOLEAN NOT NULL,
    max_account TINYINT NOT NULL,
    n_account TINYINT NOT NULL,
    is_admin BOOLEAN NOT NULL,
    allow_polling BOOLEAN NOT NULL,
    current_iksm iksm VARCHAR(32) NOT NULL,
    PRIMARY KEY (uid)
);