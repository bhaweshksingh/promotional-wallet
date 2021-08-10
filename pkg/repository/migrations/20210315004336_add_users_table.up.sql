create extension if not exists "pgcrypto";

create table if not exists users (
     id uuid primary key default gen_random_uuid(),

     user_id varchar(100) not null unique ,
     first_name varchar(100) not null,
     last_name varchar(100) not null,
     password_hash varchar(200) not null,
     password_salt bytea not null,

     created_at timestamp without time zone default (now() at time zone 'utc'),
     updated_at timestamp without time zone default (now() at time zone 'utc')
);
