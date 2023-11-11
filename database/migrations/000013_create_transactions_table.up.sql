create table transactions
(
    id          text not null
        primary key,
    campaign_id text
        constraint fk_campaigns_transactions
            references campaigns,
    user_id     text
        constraint fk_users_transactions
            references users,
    amount      bigint,
    status      text,
    code        text,
    created_at  timestamp with time zone,
    updated_at  timestamp with time zone,
    deleted_at  timestamp with time zone
)