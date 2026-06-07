-- +goose Up
SELECT 'up SQL query';
create table if not exists booking (
                                       id uuid primary key,
                                       user_id uuid not null,
                                       listing_id uuid not null,
                                       date_from timestamp not null,
                                       date_to timestamp not null,
                                       status text not null,
                                       expires_at timestamp null,
                                       created_at timestamp not null default now(),
                                       updated_at timestamp not null default now(),
                                       total_price numeric(12, 2) not null default 0,
                                       constraint booking_date_check check (date_from < date_to),
                                       constraint booking_total_price_check check (total_price >= 0)
);

create table if not exists outbox(
                                     id uuid primary key,
                                     aggregate_id uuid,
                                     event_type text not null,
                                     payload jsonb not null,
                                     created_at timestamp not null default now(),
                                     processed_at timestamp null
);

-- +goose Down
SELECT 'down SQL query';
drop table if exists outbox;
drop table if exists booking;
