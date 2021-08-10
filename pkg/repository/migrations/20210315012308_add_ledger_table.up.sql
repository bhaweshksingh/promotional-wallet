create table if not exists ledger
(
    id uuid primary key default gen_random_uuid(),
    account_id uuid  not null,
    amount     INTEGER      not null,
    activity   varchar(36)  not null,
    type       varchar(36)  not null,
    priority   integer  not null,
    expiry     timestamp without time zone,
    created_at timestamp without time zone default (now() at time zone 'utc')
)