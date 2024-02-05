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
INSERT INTO warehouses (id, name, accessible)
    VALUES (1, 'A first warehouse', true), (2, 'Nsk-House', true),
    (3, 'No more warehouses', false), (4, 'Msk-House', true),
    (6, 'Ekb-House', false), (7, 'Rus-House', true);

INSERT INTO items (code, name, size)
    VALUES (1, 'Fake Plastic Tree', 2), (2, 'Big Clay Burdon', 4),
    (3, 'A metal arm', 6), (4, 'Black Flag T-Shirt', 2),
    (6, 'Big cyan jeans', 3), (7, 'My-Own-Summer pants', 1);

INSERT INTO warehouses_items (item_code, warehouse_id, quantity)
    VALUES (1, 1, 4), (2, 1, 5), (3, 1, 2), (4, 1, 7), (1, 2, 3),
    (2, 2, 5), (7, 2, 4), (1, 3, 10), (2, 3, 8), (6, 3, 6), (1, 4, 5), 
    (2, 4, 1), (4, 4, 5), (1, 6, 10), (2, 6, 7), (3, 6, 3), (4, 6, 2), 
    (1, 7, 4), (2, 7, 6), (3, 7, 3), (4, 7, 9), (6, 7, 1), (7, 7, 1);
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE reservations;
DROP TABLE warehouses_items;
DROP TABLE warehouses;
DROP TABLE items;
-- +goose StatementEnd
