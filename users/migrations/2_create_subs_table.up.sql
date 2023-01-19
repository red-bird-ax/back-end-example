create table if not exists subscriptions (
    subscriber_id uuid not null,
    user_id       uuid not null,
    
    foreign key (subscriber_id) references users (id),
    foreign key (user_id)       references users (id)
);