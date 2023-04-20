-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.collections
(
    id                  uuid  NOT NULL,
    user_id             uuid  NOT NULL,
    contract_address    text  NOT NULL,
    created_at          int8  NULL,
    updated_at          int8  NULL,
    CONSTRAINT collections_pkey PRIMARY KEY (id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.collections;
-- +goose StatementEnd
