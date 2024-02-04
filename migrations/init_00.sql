CREATE TABLE items(
    code int PRIMARY KEY,
    name string,
    size int
);

CREATE TABLE warehouses(
    id int PRIMARY KEY,
    name string,
    accessible bool
);

CREATE TABLE warehouses_items(
    item_code int REFERENCES items(code),
    warehouse_id int REFERENCES warehouses(id),
    quantity int,
    PRIMARY KEY (item_code, warehouse_id)
);