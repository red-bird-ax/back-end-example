create table if not exists users (
    id            uuid          not null,
    user_name     varchar(100)  not null, -- make it unique (from code too)
    full_name     varchar(255)  not null,
    password_hash varchar(1024) not null, -- adjust lenght
    status_text   varchar(255)  not null,
    timestamp     timestamp     not null,

    -- avatar_url    varchar(255),
    -- cover_url     varchar(255),

    primary key (id)
);