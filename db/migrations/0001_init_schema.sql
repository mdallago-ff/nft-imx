-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.users
(
    id                uuid  NOT NULL,
    mail             text  NOT NULL,
    api_key             text  NOT NULL,
    created_at        int8 NULL,
    updated_at        int8 NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.users;
-- +goose StatementEnd
