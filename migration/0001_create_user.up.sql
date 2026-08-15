CREATE TABLE IF NOT EXISTS users (
    id            varchar(36),
    display_name  varchar(255)    NOT NULL,
    username      varchar(255)    NOT NULL,
    password      varchar(2048)   NOT NULL,
    email         varchar(2048)   NOT NULL,
    created_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at    TIMESTAMP WITH TIME ZONE,

    CONSTRAINT users_pk PRIMARY KEY (id),
    CONSTRAINT users_username_unique UNIQUE (username),
    CONSTRAINT users_email_unique UNIQUE (email)
);