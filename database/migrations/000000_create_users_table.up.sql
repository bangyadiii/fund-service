create table users
(
    id                varchar(255) not null
        primary key,
    name              varchar(255),
    occupation        varchar(255),
    email             varchar(255),
    password          varchar(255),
    avatar            varchar(255),
    role              enum ('user', 'admin', 'super-admin') default 'user' not null,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    is_google_account boolean
);

create unique index idx_email
    on users (email)