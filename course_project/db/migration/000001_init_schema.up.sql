create table applications (
    id bigserial primary key,
    ticker text not null,
    cost float not null,
    size integer not null,
    created_at timestamp default now(),
    type text not null
);
