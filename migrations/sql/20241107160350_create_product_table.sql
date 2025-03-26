-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    shop_id UUID NOT NULL,
    category_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    image_url TEXT,
    stock INT DEFAULT 0 NOT NULL,
    price DECIMAL(19, 4) DEFAULT 0.0 NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP,

    FOREIGN KEY (shop_id) REFERENCES shops(id),
    FOREIGN KEY (category_id) REFERENCES product_categories(id)
);

-- seed data, select a shop_id and category_id from the shops and product_categories tables randomly
INSERT INTO products (shop_id, category_id, name, description, image_url, stock, price)
VALUES
    (
        (SELECT id FROM shops ORDER BY random() LIMIT 1),
        (SELECT id FROM product_categories ORDER BY random() LIMIT 1),
        'Product 1',
        'Product 1 description',
        'https://via.placeholder.com/150',
        1000,
        100.00
    ),
    (
        (SELECT id FROM shops ORDER BY random() LIMIT 1),
        (SELECT id FROM product_categories ORDER BY random() LIMIT 1),
        'Product 2',
        'Product 2 description',
        'https://via.placeholder.com/150',
        2000,
        200.00
    ),
    (
        (SELECT id FROM shops ORDER BY random() LIMIT 1),
        (SELECT id FROM product_categories ORDER BY random() LIMIT 1),
        'Product 3',
        'Product 3 description',
        'https://via.placeholder.com/150',
        3000,
        300.00
    );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS products;
-- +goose StatementEnd
