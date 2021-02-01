-- +goose Up
-- +goose StatementBegin
create table slot_banners
(
    slot_id   varchar not null,
    banner_id varchar not null,
    constraint banners_pk
        primary key (slot_id, banner_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table slot_banners
-- +goose StatementEnd
