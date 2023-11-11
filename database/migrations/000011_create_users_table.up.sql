create table users
(
    id                text not null
        primary key,
    name              text,
    occupation        text,
    email             text,
    password          text,
    avatar            text,
    role              text,
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    is_google_account boolean
);

create unique index idx_email
    on users (email);