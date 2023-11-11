create table campaigns
(
    id                varchar(255) not null
        primary key,
    user_id           varchar(255)
        constraint fk_campaigns_user
            references users,
    name              text,
    short_description text,
    description       text,
    perks             text,
    backer_count      bigint,
    goal_amount       bigint,
    current_amount    bigint,
    slug              varchar(256),
    created_at        timestamp with time zone,
    updated_at        timestamp with time zone,
    deleted_at        timestamp with time zone
);

create index idx_campaigns_deleted_at
    on campaigns (deleted_at);

create unique index idx_campaigns_slug
    on campaigns (slug);

