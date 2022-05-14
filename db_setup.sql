create table if not exists token
(
    id    int auto_increment
        primary key,
    email text        not null,
    phone int         not null,
    token varchar(40) not null,
    constraint token_id_uindex
        unique (id),
    constraint token_token_uindex
        unique (token)
);

create table if not exists message
(
    id         int auto_increment
        primary key,
    content    text                               not null,
    topic      varchar(200)                       not null,
    created_at datetime default CURRENT_TIMESTAMP not null,
    creator    int                                not null,
    constraint message_id_uindex
        unique (id),
    constraint message_token_id_fk
        foreign key (creator) references token (id)
);

create index message_topic_index
    on message (topic);

create definer = tirionadmin@`%` event if not exists delete_old_messages on schedule
    every '15' MINUTE
        starts '2022-05-14 12:30:00'
    enable
    do
    DELETE FROM message WHERE (created_at + INTERVAL 24 HOUR) < NOW();

