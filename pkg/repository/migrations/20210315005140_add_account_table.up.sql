create table if not exists accounts (
    id uuid primary key default gen_random_uuid(),

    user_id uuid not null unique ,
    balance integer not null ,
    created_at timestamp without time zone default (now() at time zone 'utc'),
    updated_at timestamp without time zone default (now() at time zone 'utc')
);
