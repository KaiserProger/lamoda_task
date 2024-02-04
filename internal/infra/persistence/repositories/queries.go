package repositories

const (
	getStoredItemsQuery = `
		SELECT 
			warehouses_items.item_code,
			warehouses_items.warehouse_id,
			warehouses_items.quantity,
			warehouses.accessibility
		FROM
			warehouses_items
		LEFT JOIN
			warehouses ON warehouses_items.warehouse_id = warehouses.id 
		WHERE
			warehouse.accessible = true
			AND warehouses_items.item_code IN (?)
			AND warehouses_items.quantity > 0
		GROUP BY warehouses_items.item_code;`
	getWarehousesQuery   = `SELECT * FROM warehouses WHERE id IN (?);`
	makeReservationQuery = `
		INSERT INTO
			reservations (item_code, warehouse_id, quantity)
		VALUES (?)
		ON CONFLICT (item_code, warehouse_id) DO
		UPDATE SET quantity=EXCLUDED.quantity;`
	dereserveQuery   = `UPDATE reservations SET quantity=quantity-? WHERE item_code IN (?) AND quantity > ? RETURNING warehouse_id;`
	updateStockQuery = `UPDATE warehouses_items SET quantity=quantity+? WHERE (item_code, warehouse_id) IN (?);`
)
