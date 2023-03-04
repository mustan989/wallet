create table wallets
(
    id          bigserial primary key,
    "name"      varchar(50)    not null,
    description varchar(300),
    currency    char(3)        not null,
    amount      decimal(19, 2) not null default 0,
    personal    boolean        not null default false,
    created_at  timestamptz    not null default now(),
    updated_at  timestamptz    not null default now(),
    deleted_at  timestamptz
)