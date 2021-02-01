-- +goose Up
-- +goose StatementBegin
create table counters
(
    slot_id varchar not null,
    banner_id varchar not null,
    user_group_id varchar,
    views bigint default 0 not null,
    clicks bigint default 0 not null,
    constraint counters_slot_banners_slot_id_banner_id_fk
        foreign key (slot_id, banner_id) references slot_banners
            on delete cascade
);

create unique index counters_slot_id_banner_id_user_group_id_uindex
	on counters (slot_id, banner_id, user_group_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table counters
-- +goose StatementEnd
