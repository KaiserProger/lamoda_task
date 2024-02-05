-- +goose Up
-- +goose StatementBegin
CREATE TABLE items(
    code int PRIMARY KEY,
    name varchar(255),
    size int
);

CREATE TABLE warehouses(
    id int PRIMARY KEY,
    name varchar(255),
    accessible boolean
);

CREATE TABLE warehouses_items(
    item_code int REFERENCES items(code),
    warehouse_id int REFERENCES warehouses(id),
    quantity int,
    PRIMARY KEY (item_code, warehouse_id)
);

CREATE TABLE reservations(
    item_code int REFERENCES items(code),
    warehouse_id int REFERENCES warehouses(id),
    quantity int,
    created_at timestamp DEFAULT timezone('utc', now()),
    PRIMARY KEY (item_code, warehouse_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE reservations;
DROP TABLE warehouses_items;
DROP TABLE warehouses;
DROP TABLE items;
-- +goose StatementEnd
