create EXTENSION if not exists "uuid-ossp";
create type user_status as enum ('active','inactive','banned','pedding');

create table users (
    user_id uuid primary key default uuid_generate_v4(),
    username varchar(50) not null,
    password_hash text not null,
    firstname varchar(50) not null,
    lastname varchar(50) not null,
    phonenumber varchar(12) not null,
    email varchar(50) not null,
    role varchar(50) default 'user',
    created_at timestamp with time zone default current_timestamp,
    updated_at timestamp with time zone default current_timestamp,
    status user_status default 'active'
);