-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.tokens
(
    id                  uuid  NOT NULL,
    collection_id       uuid  NOT NULL,
    token_id            text  NOT NULL,
    created_at          int8  NULL,
    updated_at          int8  NULL,
    CONSTRAINT tokens_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.tokens;
-- +goose StatementEnd
