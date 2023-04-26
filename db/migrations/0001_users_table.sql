-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users
(
    id              uuid  NOT NULL,
    mail            text  NOT NULL,
    api_key         text  NOT NULL,
    "private"       text  NOT NULL,
    "public"        text  NOT NULL,
    address         text  NOT NULL,
    stark_key       text  NULL,
    created_at      int8  NULL,
    updated_at      int8  NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.users;
-- +goose StatementEnd
