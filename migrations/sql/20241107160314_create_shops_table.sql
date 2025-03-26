-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    terms TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now() NOT NULL,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- seed data
INSERT INTO shops (user_id, name, description, terms)
VALUES (gen_random_uuid(), 'AfterOffice Shop', 'The best shop in the world', 'The best terms in the world');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS shops;
-- +goose StatementEnd
